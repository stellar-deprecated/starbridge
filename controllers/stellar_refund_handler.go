package controllers

import (
	"context"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/stellar"
)

var (
	InvalidSequenceNumber = problem.P{
		Type:   "invalid_sequence_number",
		Title:  "Invalid Sequence Number",
		Status: http.StatusBadRequest,
		Detail: "The sequence parameter is not valid.",
	}
	InvalidSourceAccount = problem.P{
		Type:   "invalid_source_account",
		Title:  "Invalid Source Account",
		Status: http.StatusBadRequest,
		Detail: "The source account parameter is not valid.",
	}
)

type StellarRefundHandler struct {
	StellarBuilder         *stellar.Builder
	StellarSigner          *stellar.Signer
	StellarObserver        stellar.Observer
	EthereumObserver       ethereum.Observer
	WithdrawalWindow       time.Duration
	EthereumFinalityBuffer uint64
}

func (c *StellarRefundHandler) CanRefund(ctx context.Context, deposit stellar.Deposit) error {
	withdrawalDeadline := deposit.Time.Add(c.WithdrawalWindow)

	// Checks on Ethereum side:
	// - Ensure that there was no withdrawal to Ethereum account
	// - The response from the client is after the withdrawal deadline
	depositID := common.HexToHash(deposit.ID)
	requestStatus, err := c.EthereumObserver.GetRequestStatus(ctx, depositID)
	if err != nil {
		return errors.Wrap(err, "error getting request status from ethereum observer")
	}

	if requestStatus.BlockNumber <= c.EthereumFinalityBuffer {
		return EthereumNodeBehind
	}

	block, err := c.EthereumObserver.GetBlockByNumber(ctx, requestStatus.BlockNumber-c.EthereumFinalityBuffer)
	if err != nil {
		return errors.Wrap(err, "error getting block from ethereum observer")
	}

	if !block.Time.After(withdrawalDeadline) {
		return WithdrawalWindowStillActive
	}

	if requestStatus.Fulfilled {
		return WithdrawalAlreadyExecuted
	}

	return nil
}

func getSourceAccount(r *http.Request, bridgeAccount, signer string) (string, error) {
	sourceAccount := r.PostFormValue("source")
	if _, err := strkey.Decode(strkey.VersionByteAccountID, sourceAccount); err != nil {
		return "", InvalidSourceAccount
	}
	if sourceAccount == bridgeAccount || sourceAccount == signer {
		return "", InvalidSourceAccount
	}
	return sourceAccount, nil
}

func (c *StellarRefundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sequence, err := strconv.ParseInt(r.PostFormValue("sequence"), 10, 64)
	if err != nil || sequence < 0 {
		problem.Render(r.Context(), w, InvalidSequenceNumber)
		return
	}

	sourceAccount, err := getSourceAccount(r, c.StellarBuilder.BridgeAccount, c.StellarSigner.Signer.Address())
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	deposit, err := getStellarDeposit(c.StellarObserver, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	if err = c.CanRefund(r.Context(), deposit); err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		problem.Render(r.Context(), w, errors.Wrap(err, "error decoding deposit id"))
		return
	}

	tx, err := c.StellarBuilder.BuildTransaction(
		deposit.Token,
		sourceAccount,
		deposit.Sender,
		deposit.Amount,
		sequence,
		txnbuild.TimeoutInfinite,
		depositIDBytes,
	)
	if err != nil {
		problem.Render(r.Context(), w, errors.Wrap(err, "error building outgoing stellar transaction"))
		return
	}

	signature, err := c.StellarSigner.Sign(tx)
	if err != nil {
		problem.Render(r.Context(), w, errors.Wrap(err, "error signing outgoing stellar transaction"))
		return
	}

	sigs := tx.Signatures()
	tx.V1.Signatures = append(sigs, signature)

	txBase64, err := xdr.MarshalBase64(tx)
	if err != nil {
		problem.Render(r.Context(), w, errors.Wrap(err, "error marshaling outgoing stellar transaction"))
		return
	}

	_, _ = w.Write([]byte(txBase64))
}
