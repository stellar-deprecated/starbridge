package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/stellar/go/support/render/problem"

	"github.com/stellar/starbridge/stellar"
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
	InvalidEventIndex = problem.P{
		Type:   "invalid_event_index",
		Title:  "Invalid Event Index",
		Status: http.StatusBadRequest,
		Detail: "The given event index for the Stellar deposit is invalid.",
	}
	InvalidOperationIndex = problem.P{
		Type:   "invalid_operation_index",
		Title:  "Invalid Operation Index",
		Status: http.StatusBadRequest,
		Detail: "The given operation index for the Stellar deposit is invalid.",
	}
)

func getStellarDeposit(observer stellar.Observer, r *http.Request) (stellar.Deposit, error) {
	txHash := strings.TrimPrefix(r.PostFormValue("transaction_hash"), "0x")
	if !validTxHash.MatchString(txHash) {
		return stellar.Deposit{}, InvalidStellarTxHash
	}

	parsed, err := strconv.ParseInt(r.PostFormValue("event_index"), 10, 32)
	if err != nil {
		return stellar.Deposit{}, InvalidEventIndex
	}
	eventIndex := uint(parsed)

	parsed, err = strconv.ParseInt(r.PostFormValue("operation_index"), 10, 32)
	if err != nil {
		return stellar.Deposit{}, InvalidOperationIndex
	}
	operationIndex := uint(parsed)

	deposit, err := observer.GetDeposit(r.Context(), txHash, operationIndex, eventIndex)
	if stellar.IsInvalidGetDepositRequest(err) {
		return stellar.Deposit{}, InvalidDepositLog
	} else if err == stellar.ErrTxHashNotFound {
		return stellar.Deposit{}, StellarTxHashNotFound
	} else if err != nil {
		return stellar.Deposit{}, err
	}

	return deposit, nil
}
