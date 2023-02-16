package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/backend"
	"github.com/stellar/starbridge/store"
	"net/http"
)

type StellarConcordiumWithdrawalHandler struct {
	Store                         *store.DB
	ConcordiumWithdrawalValidator backend.ConcordiumWithdrawalValidator
}

type ConcordiumSignatureResponse struct {
	Address    []byte `json:"address"`
	Signature  string `json:"signature"`
	DepositID  string `json:"deposit_id"`
	Expiration int64  `json:"expiration,string"`
	Token      string `json:"token"`
	Amount     string `json:"amount"`
}

func (c *StellarConcordiumWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getStellarDeposit(c.Store, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	// Check if outgoing transaction exists
	row, err := c.Store.GetConcordiumSignature(r.Context(), store.Withdraw, deposit.ID)
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
			Amount:     row.Amount,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			_, _ = w.Write(responseBytes)
		}
		return
	}

	_, err = c.ConcordiumWithdrawalValidator.CanWithdraw(r.Context(), deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	err = c.Store.InsertSignatureRequest(r.Context(), store.SignatureRequest{
		WithdrawChain: store.Concordium,
		DepositChain:  store.Stellar,
		Action:        store.Withdraw,
		DepositID:     deposit.ID,
	})
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
