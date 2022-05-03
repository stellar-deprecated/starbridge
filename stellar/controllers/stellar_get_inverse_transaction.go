package controllers

import (
	"database/sql"
	"net/http"

	"github.com/stellar/starbridge/store"
)

type StellarGetInverseTransactionForEthereum struct {
	Store *store.Memory
}

func (c *StellarGetInverseTransactionForEthereum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ethereumTransactionHash := r.PostForm.Get("transaction_hash")

	// Ensure incoming transaction exists
	_, err := c.Store.GetIncomingEthereumTransactionByHash(ethereumTransactionHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Check if outgoing transaction exists
	outgoingTransaction, err := c.Store.GetOutgoingStellarTransactionForEthereumByHash(ethereumTransactionHash)
	if err != nil && err != sql.ErrNoRows {
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
	err = c.Store.InsertSignatureRequestForIncomingEthereumTransaction(ethereumTransactionHash)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusAccepted)
}
