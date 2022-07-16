package controllers

import (
	"database/sql"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

type StellarWithdrawalHandler struct {
	StellarClient              *horizonclient.Client
	Observer                   ethereum.Observer
	Store                      *store.DB
	StellarWithdrawalValidator backend.StellarWithdrawalValidator
	EthereumFinalityBuffer     uint64
}

func (c *StellarWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getEthereumDeposit(c.Observer, c.Store, c.EthereumFinalityBuffer, r)
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

	_, err = c.StellarWithdrawalValidator.CanWithdraw(r.Context(), deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Outgoing Stellar transaction does not exist so create signature request.
	// Duplicate requests for the same signatures are not allowed but the error is ignored.
	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		DepositChain: store.Ethereum,
		Action:       store.Withdraw,
		DepositID:    deposit.ID,
	})
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
