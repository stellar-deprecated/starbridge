package backend

import (
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"
	"time"

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

	StellarClient   *horizonclient.Client
	StellarBuilder  *txbuilder.Builder
	StellarSigner   *signer.Signer
	StellarObserver *txobserver.Observer

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
			switch sr.IncomingType {
			case store.Ethereum:
				err := w.processIncomingEthereumSignatureRequest(sr)
				if err != nil {
					w.log.WithFields(log.F{"err": err, "hash": sr.IncomingTransactionHash}).
						Error("Cannot process signature request")
				}

				w.log.WithField("hash", sr.IncomingTransactionHash).
					WithField("network", sr.IncomingType).
					Info("Processed signature request successfully")

				err = w.Store.DeleteSignatureRequestForIncomingEthereumTransaction(context.TODO(), sr.IncomingTransactionHash)
				if err != nil {
					w.log.WithFields(log.F{"err": err, "hash": sr.IncomingTransactionHash}).
						Error("Error removing signature request")
				}
			}
		}
	}
}

func (w *Worker) processIncomingEthereumSignatureRequest(sr store.SignatureRequest) error {
	hash := sr.IncomingTransactionHash

	incomingEthereumTransaction, err := w.Store.GetIncomingEthereumTransactionByHash(context.TODO(), hash)
	if err != nil {
		return errors.Wrap(err, "error getting incoming ethereum transaction")
	}

	// Ensure incoming tx can still be withdrawn
	lastLedgerCloseTime, err := w.Store.GetLastLedgerCloseTime(context.TODO())
	if err != nil {
		return errors.Wrap(err, "error getting last ledger close time")
	}

	if lastLedgerCloseTime.After(incomingEthereumTransaction.WithdrawExpiration) {
		return errors.New("withdrawal no longer possible")
	}

	outgoingStellarTransaction, err := w.Store.GetOutgoingStellarTransactionForEthereumByHash(context.TODO(), hash)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "error getting outgoing stellar transaction")
	}

	// Ensure outgoing tx is not pending or success
	if outgoingStellarTransaction.State == store.PendingState ||
		outgoingStellarTransaction.State == store.SuccessState {
		return errors.Errorf("outgoing transaction is in `%s` state", outgoingStellarTransaction.State)
	}

	// Load source account sequence
	sourceAccount, err := w.StellarClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: incomingEthereumTransaction.StellarAddress,
	})
	if err != nil {
		return errors.Wrap(err, "error getting account details")
	}

	// Ensure source account sequence was not bumped in a meantime
	lastLedgerSequence, err := w.Store.GetLastLedgerSequence(context.Background())
	if err != nil {
		return errors.Wrap(err, "error getting last ledger sequence")
	}

	if lastLedgerSequence < sourceAccount.LastModifiedLedger {
		return errors.New("skipping, account sequence possibly bumped after last ledger ingested")
	}

	// All good: build, sign and persist outgoing transaction
	amountRat, ok := new(big.Rat).SetString(incomingEthereumTransaction.ValueWei)
	if !ok {
		return errors.Errorf("cannot convert value in wei to bit.Rat: %s", incomingEthereumTransaction.ValueWei)
	}
	amountRat.Quo(amountRat, weiInEth)

	incomingEthereumTransactionHashBytes, err := hex.DecodeString(incomingEthereumTransaction.Hash)
	if err != nil {
		return errors.Wrap(err, "error decoding incomingEthereumTransaction.Hash")
	}

	tx, err := w.StellarBuilder.BuildTransaction(
		incomingEthereumTransaction.StellarAddress,
		incomingEthereumTransaction.StellarAddress,
		amountRat.FloatString(7),
		sourceAccount.Sequence+1,
		// TODO: ensure using WithdrawExpiration without any time buffer is safe
		incomingEthereumTransaction.WithdrawExpiration.Unix(),
		incomingEthereumTransactionHashBytes,
	)
	if err != nil {
		return errors.Wrap(err, "error building outgoing stellar transaction")
	}

	outgoingHash, signature, err := w.StellarSigner.Sign(tx)
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
		State:    store.PendingState,
		Hash:     outgoingHash,
		Envelope: txBase64,

		IncomingType:            sr.IncomingType,
		IncomingTransactionHash: sr.IncomingTransactionHash,
	}

	err = w.Store.UpsertOutgoingStellarTransaction(context.TODO(), outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}
