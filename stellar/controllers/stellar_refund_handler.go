package controllers

import (
	"database/sql"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/store"
)

type StellarRefundHandler struct {
	StellarClient          *horizonclient.Client
	Store                  *store.DB
	StellarRefundValidator backend.StellarRefundValidator
}

func (c *StellarRefundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getStellarDeposit(c.Store, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Check if outgoing transaction exists
	outgoingTransaction, err := c.Store.GetOutgoingStellarTransaction(r.Context(), store.Refund, deposit.ID)
	if err != nil && err != sql.ErrNoRows {
		problem.Render(r.Context(), w, err)
		return
	}
	if err == nil {
		sourceAccount, err := c.StellarClient.AccountDetail(horizonclient.AccountRequest{
			AccountID: deposit.Destination,
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

	if _, err = c.StellarRefundValidator.CanRefund(r.Context(), deposit); err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		DepositChain: store.Stellar,
		Action:       store.Refund,
		DepositID:    deposit.ID,
	})
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
