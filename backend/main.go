package backend

import (
	"database/sql"
	"math/big"
	"time"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
	"github.com/stellar/starbridge/stellar/signer"
	"github.com/stellar/starbridge/stellar/txbuilder"
	"github.com/stellar/starbridge/store"
)

var (
	ten      = big.NewInt(10)
	eighteen = big.NewInt(18)
	// weiInEth = 10^18
	weiInEth = new(big.Rat).SetInt(new(big.Int).Exp(ten, eighteen, nil))
)

type Worker struct {
	Store *store.Memory

	StellarBuilder *txbuilder.Builder
	StellarSigner  *signer.Signer

	log *log.Entry
}

func (w *Worker) Run() error {
	w.log = log.WithField("service", "backend")

	w.log.Info("Starting worker")

	for {
		signatureRequests, err := w.Store.GetSignatureRequests()
		if err != nil {
			if err == sql.ErrNoRows {
				time.Sleep(time.Second)
				continue
			} else {
				w.log.WithField("err", err).Error("cannot get signature requests")
			}
		}

		w.log.Infof("Processing %d signature requests", len(signatureRequests))

		for _, sr := range signatureRequests {
			switch sr.IncomingType {
			case store.Ethereum:
				err := w.processIncomingEthereumSignatureRequest(sr)
				if err != nil {
					w.log.WithFields(log.F{"err": err, "hash": *sr.IncomingEthereumTransactionHash}).
						Error("Cannot process signature request")
				}

				w.log.WithField("hash", *sr.IncomingEthereumTransactionHash).
					WithField("network", sr.IncomingType).
					Info("Processed signature request successfully")

				err = w.Store.DeleteSignatureRequestForIncomingEthereumTransaction(*sr.IncomingEthereumTransactionHash)
				if err != nil {
					w.log.WithFields(log.F{"err": err, "hash": *sr.IncomingEthereumTransactionHash}).
						Error("Error removing signature request")
				}
			}
		}
	}
}

func (w *Worker) processIncomingEthereumSignatureRequest(sr store.SignatureRequest) error {
	hash := *sr.IncomingEthereumTransactionHash

	incomingEthereumTransaction, err := w.Store.GetIncomingEthereumTransactionByHash(hash)
	if err != nil {
		return errors.Wrap(err, "error getting incoming ethereum transaction")
	}

	// Ensure incoming tx can still be withdrawn
	if incomingEthereumTransaction.WithdrawExpiration.After(time.Now()) {
		return errors.New("transaction withdraw time expired")
	}

	outgoingStellarTransaction, err := w.Store.GetOutgoingStellarTransactionForEthereumByHash(hash)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "error getting outgoing stellar transaction")
	}

	// Ensure outgoing tx is not pending or success
	if outgoingStellarTransaction.State == store.PendingState ||
		outgoingStellarTransaction.State == store.SuccessState {
		return errors.Errorf("outgoing transaction is in `%s` state", outgoingStellarTransaction.State)
	}

	// All good: build, sign and persist outgoing transaction
	amountRat := new(big.Rat).SetInt(incomingEthereumTransaction.ValueWei)
	amountRat.Quo(amountRat, weiInEth)

	tx, err := w.StellarBuilder.BuildTransaction(
		incomingEthereumTransaction.StellarAddress,
		incomingEthereumTransaction.StellarAddress,
		amountRat.FloatString(7),
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
		// Overflow not possible because MaxTime is set by Starbridge
		Expiration: time.Unix(int64(tx.V1.Tx.TimeBounds.MaxTime), 0),

		IncomingType:                    sr.IncomingType,
		IncomingEthereumTransactionHash: sr.IncomingEthereumTransactionHash,
	}

	err = w.Store.UpsertOutgoingStellarTransaction(outgoingTx)
	if err != nil {
		return errors.Wrap(err, "error upserting outgoing stellar transaction")
	}

	return nil
}
