package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"

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
	StellarTxNotSuccessful = problem.P{
		Type:   "stellar_tx_not_successful",
		Title:  "Stellar Transaction Not Successful",
		Status: http.StatusUnprocessableEntity,
		Detail: "The stellar transaction was not successful.",
	}
	StellarTxHasMultipleOperations = problem.P{
		Type:   "stellar_tx_has_multiple_operations",
		Title:  "Stellar Transaction Has Multiple Operations",
		Status: http.StatusUnprocessableEntity,
		Detail: "The stellar deposit transaction is invalid because it has multiple operations.",
	}
	StellarTxIsNotPayment = problem.P{
		Type:   "stellar_tx_is_not_payment",
		Title:  "Stellar Transaction Is Not Payment",
		Status: http.StatusUnprocessableEntity,
		Detail: "The stellar deposit transaction is invalid because it is not a payment.",
	}
	StellarTxIsNotPaymentToBridge = problem.P{
		Type:   "stellar_tx_is_not_payment_to_bridge",
		Title:  "Stellar Transaction Is Not Payment To Bridge",
		Status: http.StatusUnprocessableEntity,
		Detail: "The stellar deposit transaction is invalid because it " +
			"is not a payment to the bridge account.",
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

func getStellarDeposit(client *horizonclient.Client, bridgeAccount string, depositStore *store.DB, r *http.Request) (store.StellarDeposit, error) {
	txHash := strings.TrimPrefix(r.PostFormValue("transaction_hash"), "0x")
	if !validTxHash.MatchString(txHash) {
		return store.StellarDeposit{}, InvalidStellarTxHash
	}

	deposit, err := depositStore.GetStellarDeposit(r.Context(), txHash)
	if err == nil {
		return deposit, nil
	} else if err != sql.ErrNoRows {
		return store.StellarDeposit{}, err
	}

	ops, err := client.Operations(horizonclient.OperationRequest{
		ForTransaction: txHash,
		IncludeFailed:  false,
		Limit:          200,
		Join:           "transactions",
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok && herr.Response.StatusCode == http.StatusNotFound {
			return store.StellarDeposit{}, StellarTxHashNotFound
		}
		return store.StellarDeposit{}, err
	}

	if len(ops.Embedded.Records) != 1 {
		return store.StellarDeposit{}, StellarTxHasMultipleOperations
	}

	operation := ops.Embedded.Records[0]

	payment, ok := operation.(operations.Payment)
	if !ok {
		return store.StellarDeposit{}, StellarTxIsNotPayment
	}

	if payment.To != bridgeAccount {
		return store.StellarDeposit{}, StellarTxIsNotPaymentToBridge
	}

	if !payment.TransactionSuccessful {
		return store.StellarDeposit{}, StellarTxNotSuccessful
	}

	var assetString string
	if payment.Asset.Type == "native" {
		assetString = "native"
	} else {
		assetString = payment.Asset.Code + ":" + payment.Asset.Issuer
	}

	storeDeposit := store.StellarDeposit{
		ID:          txHash,
		Asset:       assetString,
		LedgerTime:  payment.LedgerCloseTime.Unix(),
		Sender:      payment.From,
		Destination: payment.Transaction.Memo,
		Amount:      payment.Amount,
	}
	if err = depositStore.InsertStellarDeposit(r.Context(), storeDeposit); err != nil {
		return store.StellarDeposit{}, err
	}

	return storeDeposit, nil
}
