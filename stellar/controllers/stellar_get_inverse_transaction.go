package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

// TODO remove after prototype demo
type TestDeposit struct {
	Store *store.DB
}

func (c *TestDeposit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := r.PostFormValue("hash")
	stellarAddress := r.PostFormValue("stellar_address")

	incomingTx := store.IncomingEthereumTransaction{
		Hash:           hash,
		ValueWei:       1000,
		StellarAddress: stellarAddress,
	}

	err := c.Store.InsertIncomingEthereumTransaction(r.Context(), incomingTx)
	if err != nil {
		log.WithField("error", err).Error("Error inserting incoming ethereum transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(incomingTx.Hash))
}

type StellarGetInverseTransactionForEthereum struct {
	Store *store.DB
}

func (c *StellarGetInverseTransactionForEthereum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ethereumTransactionHash := r.PostFormValue("transaction_hash")
	txExpirationTimestamp := r.PostFormValue("tx_expiration_timestamp")
	txExpirationTimestampInt, err := strconv.ParseInt(txExpirationTimestamp, 10, 64)
	if err != nil {
		log.WithField("error", err).Error("error parsing tx_expiration_timestamp")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	txExpirationTime := time.Unix(txExpirationTimestampInt, 0)

	// Ensure incoming transaction exists
	_, err = c.Store.GetIncomingEthereumTransactionByHash(r.Context(), ethereumTransactionHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.WithField("error", err).Error("Error getting an incoming ethereum transaction")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Check TxExpirationTimestamp
	// TODO - replace with last ledger close time persisted in a DB
	lastLedgerCloseTime := time.Now().UTC()

	if txExpirationTime.Before(lastLedgerCloseTime) {
		log.Error("tx expiration timestamp in the past")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO 10m below should be configurable
	if txExpirationTime.After(lastLedgerCloseTime.Add(10 * time.Minute)) {
		log.Error("tx expiration timestamp too far in the future")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if outgoing transaction exists
	outgoingTransaction, err := c.Store.GetOutgoingStellarTransactionForEthereumByHash(r.Context(), ethereumTransactionHash)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("error", err).Error("Error getting an outgoing stellar transaction for ethereum")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return outgoing transaction if pending or success, otherwise create a signature request
	if err == nil &&
		(outgoingTransaction.State == store.PendingState || outgoingTransaction.State == store.SuccessState) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(outgoingTransaction.Envelope))
		return
	}

	// Check TxExpirationTimestamp
	// TODO - replace with last ledger close time persisted in a DB
	lastLedgerCloseTime = time.Now().UTC()

	if txExpirationTime.Before(lastLedgerCloseTime) {
		log.Error("tx expiration timestamp in the past")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO 10m below should be configurable
	if txExpirationTime.After(lastLedgerCloseTime.Add(10 * time.Minute)) {
		log.Error("tx expiration timestamp too far in the future")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Outgoing Stellar transaction does not exist so create signature request.
	// Duplicate requests for the same signatures are not allowed but the error is ignored.
	err = c.Store.InsertSignatureRequestForIncomingEthereumTransaction(r.Context(), ethereumTransactionHash, txExpirationTimestampInt)
	if err != nil {
		// Ignore duplicate violations
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.WithField("error", err).Error("Error inserting a signature request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
