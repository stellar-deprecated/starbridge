package controllers

import (
	"database/sql"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/concordium"
	"math/big"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/store"
)

type ConcordiumStellarWithdrawalHandler struct {
	StellarClient                          *horizonclient.Client
	Observer                               concordium.Observer
	Store                                  *store.DB
	ConcordiumToStellarWithdrawalValidator backend.StellarWithdrawalValidator
}

func (c *ConcordiumStellarWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txHash := r.PostFormValue("transaction_hash")
	stellarAddress := r.PostFormValue("destination")
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, stellarAddress)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	var intEncoded big.Int
	intEncoded.SetBytes(decoded)
	deposit, err := getConcordiumDeposit(txHash, c.Store)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Check if outgoing transaction exists
	outgoingTransaction, err := c.Store.GetOutgoingStellarTransaction(r.Context(), store.Withdraw, deposit.ID)
	if err != nil && err != sql.ErrNoRows {
		problem.Render(r.Context(), w, err)
		return
	}
	if err == nil {
		sourceAccount, err := c.StellarClient.AccountDetail(horizonclient.AccountRequest{
			AccountID: outgoingTransaction.SourceAccount,
		})
		if err != nil {
			problem.Render(r.Context(), w, err)
			return
		}
		if sourceAccount.Sequence < outgoingTransaction.Sequence {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(outgoingTransaction.Envelope))
			return
		}
	}

	_, err = c.ConcordiumToStellarWithdrawalValidator.CanWithdrawConcordium(r.Context(), deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Outgoing Stellar transaction does not exist so create signature request.
	// Duplicate requests for the same signatures are not allowed but the error is ignored.
	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		WithdrawChain: store.Stellar,
		DepositChain:  store.Concordium,
		Action:        store.Withdraw,
		DepositID:     deposit.ID,
	})
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
