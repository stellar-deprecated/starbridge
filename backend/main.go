package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/stellar/starbridge/concordium"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/stellar/txobserver"
	"github.com/stellar/starbridge/store"
)

type Worker struct {
	Store *store.DB

	StellarClient              *horizonclient.Client
	StellarBuilder             *txbuilder.Builder
	StellarSigner              *signer.Signer
	StellarObserver            *txobserver.Observer
	StellarWithdrawalValidator StellarWithdrawalValidator
	StellarRefundValidator     StellarRefundValidator

	EthereumRefundValidator     EthereumRefundValidator
	EthereumWithdrawalValidator EthereumWithdrawalValidator
	EthereumSigner              ethereum.Signer

	OkxWithdrawalValidator OkxWithdrawalValidator
	OkxSigner              ethereum.Signer

	ConcordiumWithdrawalValidator ConcordiumWithdrawalValidator
	ConcordiumSigner              concordium.Signer

	log *log.Entry
}

func (w *Worker) Run(ctx context.Context) {
	w.log = log.WithField("service", "backend")

	w.log.Info("Starting worker")

	for ctx.Err() == nil {
		// Process all new ledgers before processing signature requests
		w.StellarObserver.ProcessNewLedgers(ctx)

		signatureRequests, err := w.Store.GetSignatureRequests(ctx)
		if err != nil {
			w.log.WithField("err", err).Error("cannot get signature requests")
			time.Sleep(time.Second)
			continue
		}

		if len(signatureRequests) == 0 {
			time.Sleep(time.Second)
			continue
		}

		w.log.Infof("Processing %d signature requests", len(signatureRequests))

		for _, sr := range signatureRequests {
			var err error
			switch sr.Action {
			case store.Withdraw:
				switch sr.DepositChain {
				case store.Ethereum:
					err = w.processStellarWithdrawalRequest(ctx, sr)
				case store.Concordium:
					err = w.processConcordiumToStellarWithdrawalRequest(ctx, sr)
				case store.Okx:
					err = w.processOkxToStellarWithdrawalRequest(ctx, sr)
				case store.Stellar:
					switch sr.WithdrawChain {
					case store.Ethereum:
						err = w.processEthereumWithdrawalRequest(ctx, sr)
					case store.Concordium:
						err = w.processStellarToConcordiumWithdrawalRequest(ctx, sr)
					case store.Okx:
						err = w.processStellarToOkxWithdrawalRequest(ctx, sr)
					default:
						err = fmt.Errorf("withdrawals for deposit chain %v and withdrawal chain %v is not supported", sr.DepositChain, sr.WithdrawChain)
					}
				default:
					err = fmt.Errorf("withdrawals for deposit chain %v is not supported", sr.DepositChain)
				}
			case store.Refund:
				switch sr.DepositChain {
				case store.Ethereum:
					err = w.processEthereumRefundRequest(ctx, sr)
				case store.Stellar:
					err = w.processStellarRefundRequest(ctx, sr)
				default:
					err = fmt.Errorf("refunds for deposit chain %v is not supported", sr.DepositChain)
				}
			default:
				err = fmt.Errorf("action %v is not supported", sr.Action)
			}

			if err != nil {
				w.log.WithFields(log.F{"err": err, "request": sr}).
					Error("Cannot process signature request")
				if p, ok := err.(problem.P); ok && p.Status >= 400 && p.Status < 500 {
					w.deleteRequest(ctx, sr)
				}
			} else {
				w.log.WithField("request", sr).
					Info("Processed signature request successfully")
				w.deleteRequest(ctx, sr)
			}
		}
	}
}

func (w *Worker) deleteRequest(ctx context.Context, sr store.SignatureRequest) {
	err := w.Store.DeleteSignatureRequest(ctx, sr)
	if err != nil {
		w.log.WithFields(log.F{"err": err, "request": sr}).
			Error("Error removing signature request")
	}
}

func (w *Worker) processStellarWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Ethereum {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetEthereumDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	details, err := w.StellarWithdrawalValidator.CanWithdraw(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating withdraw conditions")
	}

	// Load source account sequence
	sourceAccount, err := w.StellarClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: details.Recipient,
	})
	if err != nil {
		return errors.Wrap(err, "error getting account details")
	}
	if sourceAccount.SequenceLedger > 0 {
		if details.LedgerSequence < sourceAccount.SequenceLedger {
			return errors.New("skipping, account sequence ledger is higher than last ledger ingested")
		}
	} else {
		if details.LedgerSequence < sourceAccount.LastModifiedLedger {
			return errors.New("skipping, account sequence possibly bumped after last ledger ingested")
		}
	}

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error decoding deposit id")
	}
	tx, err := w.StellarBuilder.BuildTransaction(
		details.Asset,
		details.Recipient,
		details.Recipient,
		amount.StringFromInt64(details.Amount),
		sourceAccount.Sequence+1,
		// TODO: ensure using WithdrawExpiration without any time buffer is safe
		details.Deadline.Unix(),
		depositIDBytes,
	)
	if err != nil {
		return errors.Wrap(err, "error building outgoing stellar transaction")
	}

	signature, err := w.StellarSigner.Sign(tx)
	if err != nil {
		return errors.Wrap(err, "error signing outgoing stellar transaction")
	}

	// TODO, we need xdr.TransactionEnvelope.AppendSignature.
	sigs := tx.Signatures()
	tx.V1.Signatures = append(sigs, signature)

	txBase64, err := xdr.MarshalBase64(tx)
	if err != nil {
		return errors.Wrap(err, "error marshaling outgoing stellar transaction")
	}

	outgoingTx := store.OutgoingStellarTransaction{
		Envelope:      txBase64,
		Action:        sr.Action,
		DepositID:     sr.DepositID,
		SourceAccount: details.Recipient,
		Sequence:      tx.SeqNum(),
	}
	err = w.Store.UpsertOutgoingStellarTransaction(ctx, outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}

func (w *Worker) processOkxToStellarWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Okx {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetOkxDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting okx deposit")
	}

	details, err := w.StellarWithdrawalValidator.CanWithdrawOkx(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating withdraw conditions")
	}

	// Load source account sequence
	sourceAccount, err := w.StellarClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: details.Recipient,
	})
	if err != nil {
		return errors.Wrap(err, "error getting account details")
	}
	if sourceAccount.SequenceLedger > 0 {
		if details.LedgerSequence < sourceAccount.SequenceLedger {
			return errors.New("skipping, account sequence ledger is higher than last ledger ingested")
		}
	} else {
		if details.LedgerSequence < sourceAccount.LastModifiedLedger {
			return errors.New("skipping, account sequence possibly bumped after last ledger ingested")
		}
	}

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error decoding deposit id")
	}
	tx, err := w.StellarBuilder.BuildTransaction(
		details.Asset,
		details.Recipient,
		details.Recipient,
		amount.StringFromInt64(details.Amount),
		sourceAccount.Sequence+1,
		// TODO: ensure using WithdrawExpiration without any time buffer is safe
		details.Deadline.Unix(),
		depositIDBytes,
	)
	if err != nil {
		return errors.Wrap(err, "error building outgoing stellar transaction")
	}

	signature, err := w.StellarSigner.Sign(tx)
	if err != nil {
		return errors.Wrap(err, "error signing outgoing stellar transaction")
	}

	// TODO, we need xdr.TransactionEnvelope.AppendSignature.
	sigs := tx.Signatures()
	tx.V1.Signatures = append(sigs, signature)

	txBase64, err := xdr.MarshalBase64(tx)
	if err != nil {
		return errors.Wrap(err, "error marshaling outgoing stellar transaction")
	}

	outgoingTx := store.OutgoingStellarTransaction{
		Envelope:      txBase64,
		Action:        sr.Action,
		DepositID:     sr.DepositID,
		SourceAccount: details.Recipient,
		Sequence:      tx.SeqNum(),
	}
	err = w.Store.UpsertOutgoingStellarTransaction(ctx, outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}

func (w *Worker) processConcordiumToStellarWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Concordium {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetConcordiumDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	details, err := w.StellarWithdrawalValidator.CanWithdrawConcordium(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating withdraw conditions")
	}

	// Load source account sequence
	sourceAccount, err := w.StellarClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: details.Recipient,
	})
	if err != nil {
		return errors.Wrap(err, "error getting account details")
	}
	if sourceAccount.SequenceLedger > 0 {
		if details.LedgerSequence < sourceAccount.SequenceLedger {
			return errors.New("skipping, account sequence ledger is higher than last ledger ingested")
		}
	} else {
		if details.LedgerSequence < sourceAccount.LastModifiedLedger {
			return errors.New("skipping, account sequence possibly bumped after last ledger ingested")
		}
	}

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error decoding deposit id")
	}
	tx, err := w.StellarBuilder.BuildTransaction(
		details.Asset,
		details.Recipient,
		details.Recipient,
		amount.StringFromInt64(details.Amount),
		sourceAccount.Sequence+1,
		// TODO: ensure using WithdrawExpiration without any time buffer is safe
		details.Deadline.Unix(),
		depositIDBytes,
	)
	if err != nil {
		return errors.Wrap(err, "error building outgoing stellar transaction")
	}

	signature, err := w.StellarSigner.Sign(tx)
	if err != nil {
		return errors.Wrap(err, "error signing outgoing stellar transaction")
	}

	// TODO, we need xdr.TransactionEnvelope.AppendSignature.
	sigs := tx.Signatures()
	tx.V1.Signatures = append(sigs, signature)

	txBase64, err := xdr.MarshalBase64(tx)
	if err != nil {
		return errors.Wrap(err, "error marshaling outgoing stellar transaction")
	}

	outgoingTx := store.OutgoingStellarTransaction{
		Envelope:      txBase64,
		Action:        sr.Action,
		DepositID:     sr.DepositID,
		SourceAccount: details.Recipient,
		Sequence:      tx.SeqNum(),
	}
	err = w.Store.UpsertOutgoingStellarTransaction(ctx, outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}

func (w *Worker) processStellarRefundRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Stellar {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetStellarDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	details, err := w.StellarRefundValidator.CanRefund(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating refund conditions")
	}

	// Load source account sequence
	sourceAccount, err := w.StellarClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: deposit.Sender,
	})
	if err != nil {
		return errors.Wrap(err, "error getting account details")
	}
	if sourceAccount.SequenceLedger > 0 {
		if details.LedgerSequence < sourceAccount.SequenceLedger {
			return errors.New("skipping, account sequence ledger is higher than last ledger ingested")
		}
	} else {
		if details.LedgerSequence < sourceAccount.LastModifiedLedger {
			return errors.New("skipping, account sequence possibly bumped after last ledger ingested")
		}
	}

	// All good: build, sign and persist outgoing transaction
	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error decoding deposit id")
	}

	tx, err := w.StellarBuilder.BuildTransaction(
		deposit.Asset,
		deposit.Sender,
		deposit.Sender,
		deposit.Amount,
		sourceAccount.Sequence+1,
		txnbuild.TimeoutInfinite,
		depositIDBytes,
	)
	if err != nil {
		return errors.Wrap(err, "error building outgoing stellar transaction")
	}

	signature, err := w.StellarSigner.Sign(tx)
	if err != nil {
		return errors.Wrap(err, "error signing outgoing stellar transaction")
	}

	// TODO, we need xdr.TransactionEnvelope.AppendSignature.
	sigs := tx.Signatures()
	tx.V1.Signatures = append(sigs, signature)

	txBase64, err := xdr.MarshalBase64(tx)
	if err != nil {
		return errors.Wrap(err, "error marshaling outgoing stellar transaction")
	}

	outgoingTx := store.OutgoingStellarTransaction{
		Envelope:      txBase64,
		Action:        sr.Action,
		DepositID:     sr.DepositID,
		SourceAccount: deposit.Sender,
		Sequence:      tx.SeqNum(),
	}
	err = w.Store.UpsertOutgoingStellarTransaction(ctx, outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}

func (w *Worker) processEthereumRefundRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Ethereum {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetEthereumDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	if err = w.EthereumRefundValidator.CanRefund(ctx, deposit); err != nil {
		return errors.Wrap(err, "error validating refund conditions")
	}

	refundAmount, ok := new(big.Int).SetString(deposit.Amount, 10)
	if !ok {
		return errors.Errorf("cannot convert value in wei to bit.Rat: %s", deposit.Amount)
	}

	expiration := int64(math.MaxInt64)
	sig, err := w.EthereumSigner.SignWithdrawal(
		common.HexToHash(deposit.ID),
		expiration,
		common.HexToAddress(deposit.Sender),
		common.HexToAddress(deposit.Token),
		refundAmount,
	)
	if err != nil {
		return errors.Wrap(err, "error signing refund")
	}

	err = w.Store.UpsertEthereumSignature(ctx, store.EthereumSignature{
		Address:    w.EthereumSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		Action:     sr.Action,
		DepositID:  sr.DepositID,
		Expiration: expiration,
		Token:      deposit.Token,
		Amount:     deposit.Amount,
	})
	if err != nil {
		return errors.Wrap(err, "error upserting etherum signature")
	}

	return nil
}

func (w *Worker) processEthereumWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Stellar {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetStellarDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting stellar deposit")
	}

	details, err := w.EthereumWithdrawalValidator.CanWithdraw(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating withdrawal conditions")
	}

	sig, err := w.EthereumSigner.SignWithdrawal(
		common.HexToHash(deposit.ID),
		details.Deadline.Unix(),
		details.Recipient,
		details.Token,
		details.Amount,
	)
	if err != nil {
		return errors.Wrap(err, "error signing withdrawal")
	}

	err = w.Store.UpsertEthereumSignature(ctx, store.EthereumSignature{
		Address:    w.EthereumSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		Action:     sr.Action,
		DepositID:  sr.DepositID,
		Expiration: details.Deadline.Unix(),
		Token:      details.Token.String(),
		Amount:     details.Amount.String(),
	})
	if err != nil {
		return errors.Wrap(err, "error upserting etherum signature")
	}

	return nil
}

func (w *Worker) processStellarToOkxWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Stellar {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetStellarDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting stellar deposit")
	}

	details, err := w.OkxWithdrawalValidator.CanWithdraw(ctx, deposit)
	log.Info(details)
	if err != nil {
		return errors.Wrap(err, "error validating withdrawal conditions")
	}

	sig, err := w.OkxSigner.SignWithdrawal(
		common.HexToHash(deposit.ID),
		details.Deadline.Unix(),
		details.Recipient,
		details.Token,
		details.Amount,
	)
	if err != nil {
		return errors.Wrap(err, "error signing withdrawal")
	}

	err = w.Store.UpsertOkxSignature(ctx, store.OkxSignature{
		Address:    w.OkxSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		Action:     sr.Action,
		DepositID:  sr.DepositID,
		Expiration: details.Deadline.Unix(),
		Token:      details.Token.String(),
		Amount:     details.Amount.String(),
	})
	if err != nil {
		return errors.Wrap(err, "error upserting okx signature")
	}

	return nil
}

func (w *Worker) processStellarToConcordiumWithdrawalRequest(ctx context.Context, sr store.SignatureRequest) error {
	if sr.DepositChain != store.Stellar {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetStellarDeposit(ctx, sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting stellar deposit")
	}

	details, err := w.ConcordiumWithdrawalValidator.CanWithdraw(ctx, deposit)
	if err != nil {
		return errors.Wrap(err, "error validating withdrawal conditions")
	}

	sig, err := w.ConcordiumSigner.SignWithdrawal(
		ctx,
		common.HexToHash(deposit.ID),
		uint64(details.Deadline.Unix()),
		details.Recipient,
		details.Recipient,
		int(details.Amount.Int64()),
	)
	if err != nil {
		return errors.Wrap(err, "error signing withdrawal")
	}

	err = w.Store.UpsertConcordiumSignature(ctx, store.ConcordiumSignature{
		Address:    deposit.Destination,
		Signature:  hex.EncodeToString(sig),
		Action:     sr.Action,
		DepositID:  sr.DepositID,
		Expiration: details.Deadline.Unix(),
		Token:      details.Token,
		Amount:     details.Amount.String(),
	})
	if err != nil {
		return errors.Wrap(err, "error upserting etherum signature")
	}

	return nil
}
