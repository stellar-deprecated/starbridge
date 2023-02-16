package controllers

import (
	"database/sql"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/store"
	"net/http"
	"strings"
)

var (
	InvalidStellarTxHash = problem.P{
		Type:   "invalid_stellar_tx_hash",
		Title:  "Invalid Stellar Transaction Hash",
		Status: http.StatusBadRequest,
		Detail: "The transaction hash of the Stellar transaction is invalid.",
	}
	StellarTxHashNotFound = problem.P{
		Type:   "stellar_tx_hash_not_found",
		Title:  "Stellar Transaction Hash Not Found",
		Status: http.StatusNotFound,
		Detail: "The stellar transaction cannot be found.",
	}
)

func getStellarDeposit(depositStore *store.DB, r *http.Request) (store.StellarDeposit, error) {
	txHash := strings.TrimPrefix(r.PostFormValue("transaction_hash"), "0x")
	destination := r.PostFormValue("destination")
	if !validTxHash.MatchString(txHash) {
		return store.StellarDeposit{}, InvalidStellarTxHash
	}

	deposit, err := depositStore.GetStellarDeposit(r.Context(), txHash)
	if err == sql.ErrNoRows {
		return store.StellarDeposit{}, StellarTxHashNotFound
	}
	if deposit.Destination == "" {
		err := depositStore.UpdateStellarDepositDestination(r.Context(), deposit.ID, destination)
		if err != nil {
			return store.StellarDeposit{}, err
		}
	}
	return deposit, err
}
