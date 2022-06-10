package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

// TODO remove after prototype demo
type TestDeposit struct {
	Store *store.DB
}

func (c *TestDeposit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stellarAddress := r.PostFormValue("stellar_address")

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.WithField("error", err).Error("Error generating random bytes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	incomingTx := store.IncomingEthereumTransaction{
		Hash:           hex.EncodeToString(bytes),
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

	// Ensure incoming transaction exists
	_, err := c.Store.GetIncomingEthereumTransactionByHash(r.Context(), ethereumTransactionHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.WithField("error", err).Error("Error getting an incoming ethereum transaction")
			w.WriteHeader(http.StatusInternalServerError)
		}
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

	// Outgoing Stellar transaction does not exist so create signature request.
	// Duplicate requests for the same signatures are not allowed but the error is ignored.
	err = c.Store.InsertSignatureRequestForIncomingEthereumTransaction(r.Context(), ethereumTransactionHash)
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
