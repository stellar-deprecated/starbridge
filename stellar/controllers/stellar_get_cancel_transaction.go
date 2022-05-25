package controllers

import (
	"database/sql"
	"net/http"

	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

type StellarGetCancelTransactionForEthereum struct {
	Store *store.DB
}

func (c *StellarGetCancelTransactionForEthereum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ethereumTransactionHash := r.PostFormValue("transaction_hash")

	// Ensure incoming transaction exists
	_, err := c.Store.GetIncomingEthereumTransactionByHash(r.Context(), ethereumTransactionHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.WithField("error", err).Error("Error getting an incomming ethereum transaction")
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

	// Error if outgoing tx exists and is pending or successful
	if err == nil &&
		(outgoingTransaction.State == store.PendingState || outgoingTransaction.State == store.SuccessState) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO
	w.WriteHeader(http.StatusAccepted)
}
