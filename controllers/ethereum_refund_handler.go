package controllers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/stellar/go/support/render/problem"

	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/stellar"
)

var (
	WithdrawalWindowStillActive = problem.P{
		Type:   "withdrawal_window_still_active",
		Title:  "Withdrawal Window Still Active",
		Status: http.StatusBadRequest,
		Detail: "The withdrawal window is still active." +
			" Wait until the withdrawal window has closed before attempting a refund.",
	}
	WithdrawalAlreadyExecuted = problem.P{
		Type:   "withdrawal_already_executed",
		Title:  "Withdrawal Already Executed",
		Status: http.StatusBadRequest,
		Detail: "The withdrawal has already been executed.",
	}
)

type EthereumSignatureResponse struct {
	Address    string `json:"address"`
	Signature  string `json:"signature"`
	DepositID  string `json:"deposit_id"`
	Expiration int64  `json:"expiration,string"`
	Token      string `json:"token"`
	Amount     string `json:"amount"`
}

type EthereumRefundHandler struct {
	EthereumObserver       ethereum.Observer
	StellarObserver        stellar.Observer
	EthereumSigner         ethereum.Signer
	EthereumFinalityBuffer uint64
	WithdrawalWindow       time.Duration
}

func (c *EthereumRefundHandler) CanRefund(ctx context.Context, deposit ethereum.Deposit) error {
	status, err := c.StellarObserver.GetRequestStatus(ctx, deposit.ID)
	if err != nil {
		return err
	}

	withdrawalDeadline := deposit.Time.Add(c.WithdrawalWindow)
	if !status.CloseTime.After(withdrawalDeadline) {
		return WithdrawalWindowStillActive
	}

	if status.Fulfilled {
		return WithdrawalAlreadyExecuted
	}

	return nil
}

func (c *EthereumRefundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getEthereumDeposit(c.EthereumObserver, c.EthereumFinalityBuffer, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	if err = c.CanRefund(r.Context(), deposit); err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	expiration := int64(math.MaxInt64)
	sig, err := c.EthereumSigner.SignWithdrawal(
		deposit.ID,
		expiration,
		deposit.Sender,
		deposit.Token,
		deposit.Amount,
	)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	responseBytes, err := json.Marshal(EthereumSignatureResponse{
		Address:    c.EthereumSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		DepositID:  hex.EncodeToString(deposit.ID[:]),
		Expiration: expiration,
		Token:      deposit.Token.String(),
		Amount:     deposit.Amount.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, _ = w.Write(responseBytes)
	}
}
