package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/stellar/starbridge/backend"

	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/store"
)

type EthereumWithdrawalHandler struct {
	Store                       *store.DB
	EthereumWithdrawalValidator backend.EthereumWithdrawalValidator
}

func (c *EthereumWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getStellarDeposit(c.Store, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Check if outgoing transaction exists
	row, err := c.Store.GetEthereumSignature(r.Context(), store.Withdraw, deposit.ID)
	if err != nil && err != sql.ErrNoRows {
		problem.Render(r.Context(), w, err)
		return
	}
	if err == nil {
		responseBytes, err := json.Marshal(EthereumSignatureResponse{
			Address:    row.Address,
			Signature:  row.Signature,
			DepositID:  row.DepositID,
			Expiration: row.Expiration,
			Token:      row.Token,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			_, _ = w.Write(responseBytes)
		}
		return
	}

	_, err = c.EthereumWithdrawalValidator.CanWithdraw(r.Context(), deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		DepositChain: store.Stellar,
		Action:       store.Withdraw,
		DepositID:    deposit.ID,
	})
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
