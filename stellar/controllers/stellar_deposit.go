package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/store"
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
	InvalidEthereumRecipient = problem.P{
		Type:   "invalid_ethereum_recipient",
		Title:  "Invalid Ethereum Recipient",
		Status: http.StatusUnprocessableEntity,
		Detail: "The recipient of the deposit is not a valid Ethereum address.",
	}
	// TODO: remove this once getStellarDeposit is used
	_ = getStellarDeposit
)

func getStellarDeposit(depositStore *store.DB, r *http.Request) (store.StellarDeposit, error) {
	txHash := strings.TrimPrefix(r.PostFormValue("transaction_hash"), "0x")
	if !validTxHash.MatchString(txHash) {
		return store.StellarDeposit{}, InvalidStellarTxHash
	}

	deposit, err := depositStore.GetStellarDeposit(r.Context(), txHash)
	if err == sql.ErrNoRows {
		return store.StellarDeposit{}, StellarTxHashNotFound
	}
	return deposit, err
}
