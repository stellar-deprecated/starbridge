package controllers

import (
	"database/sql"
	"net/http"

	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/store"
)

type StellarGetInverseTransactionForEthereum struct {
	Store *store.DB
}

func (c *StellarGetInverseTransactionForEthereum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ethereumTransactionHash := r.PostForm.Get("transaction_hash")

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
		log.WithField("error", err).Error("Error inserting a signature request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
