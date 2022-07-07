package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stellar/starbridge/ethereum"

	"github.com/stellar/go/support/render/problem"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/stellar/txobserver"
	"github.com/stellar/starbridge/store"
)

var (
	ten      = big.NewInt(10)
	eighteen = big.NewInt(18)
	// weiInEth = 10^18
	weiInEth = new(big.Rat).SetInt(new(big.Int).Exp(ten, eighteen, nil))
)

type Worker struct {
	Ctx context.Context

	Store *store.DB

	StellarClient              *horizonclient.Client
	StellarBuilder             *txbuilder.Builder
	StellarSigner              *signer.Signer
	StellarObserver            *txobserver.Observer
	StellarWithdrawalValidator StellarWithdrawalValidator
	EthereumRefundValidator    EthereumRefundValidator
	EthereumSigner             ethereum.Signer

	log *log.Entry
}

func (w *Worker) Run() {
	w.log = log.WithField("service", "backend")

	w.log.Info("Starting worker")

	for w.Ctx.Err() == nil {
		// Process all new ledgers before processing signature requests
		w.StellarObserver.ProcessNewLedgers()

		signatureRequests, err := w.Store.GetSignatureRequests(context.TODO())
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
					err = w.processStellarWithdrawalRequest(sr)
				default:
					err = fmt.Errorf("withdrawals for deposit chain %v is not supported", sr.DepositChain)
				}
			case store.Refund:
				switch sr.DepositChain {
				case store.Ethereum:
					err = w.processEthereumRefundRequest(sr)
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
					w.deleteRequest(sr)
				}
			} else {
				w.log.WithField("request", sr).
					Info("Processed signature request successfully")
				w.deleteRequest(sr)
			}
		}
	}
}

func (w *Worker) deleteRequest(sr store.SignatureRequest) {
	err := w.Store.DeleteSignatureRequest(context.TODO(), sr)
	if err != nil {
		w.log.WithFields(log.F{"err": err, "request": sr}).
			Error("Error removing signature request")
	}
}

func (w *Worker) processStellarWithdrawalRequest(sr store.SignatureRequest) error {
	if sr.DepositChain != store.Ethereum {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetEthereumDeposit(context.TODO(), sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	details, err := w.StellarWithdrawalValidator.CanWithdraw(context.TODO(), deposit)
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

	// All good: build, sign and persist outgoing transaction
	amountRat, ok := new(big.Rat).SetString(deposit.Amount)
	if !ok {
		return errors.Errorf("cannot convert value in wei to bit.Rat: %s", deposit.Amount)
	}
	amountRat.Quo(amountRat, weiInEth)

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		return errors.Wrap(err, "error decoding deposit id")
	}

	tx, err := w.StellarBuilder.BuildTransaction(
		details.Recipient,
		details.Recipient,
		amountRat.FloatString(7),
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
	err = w.Store.UpsertOutgoingStellarTransaction(context.TODO(), outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}

func (w *Worker) processEthereumRefundRequest(sr store.SignatureRequest) error {
	if sr.DepositChain != store.Ethereum {
		return fmt.Errorf("deposits from %v are not supported", sr.DepositChain)
	}
	deposit, err := w.Store.GetEthereumDeposit(context.TODO(), sr.DepositID)
	if err != nil {
		return errors.Wrap(err, "error getting ethereum deposit")
	}

	if err = w.EthereumRefundValidator.CanRefund(context.TODO(), deposit); err != nil {
		return errors.Wrap(err, "error validating refund conditions")
	}

	amount, ok := new(big.Int).SetString(deposit.Amount, 10)
	if !ok {
		return errors.Errorf("cannot convert value in wei to bit.Rat: %s", deposit.Amount)
	}

	expiration := int64(math.MaxInt64)
	sig, err := w.EthereumSigner.SignWithdrawal(
		common.HexToHash(deposit.ID),
		expiration,
		common.HexToAddress(deposit.Sender),
		common.HexToAddress(deposit.Token),
		amount,
	)
	if err != nil {
		return errors.Wrap(err, "error signing refund")
	}

	err = w.Store.UpsertEthereumSignature(context.TODO(), store.EthereumSignature{
		Address:    w.EthereumSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		Action:     sr.Action,
		DepositID:  sr.DepositID,
		Expiration: expiration,
	})
	if err != nil {
		return errors.Wrap(err, "error upserting etherum signature")
	}

	return nil
}
