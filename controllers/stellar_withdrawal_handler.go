package controllers

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/go/xdr"

	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/stellar"
)

var (
	InvalidStellarRecipient = problem.P{
		Type:   "invalid_stellar_recipient",
		Title:  "Invalid Stellar Recipient",
		Status: http.StatusBadRequest,
		Detail: "The recipient of the deposit is not a valid Stellar address.",
	}
)

// StellarWithdrawalDetails includes metadata about the
// validation result.
type StellarWithdrawalDetails struct {
	// Deadline is the deadline for executing the withdrawal
	// transaction on Stellar.
	Deadline time.Time
	// Recipient is the Stellar account which should receive the
	// withdrawal.
	Recipient string
	// AssetContractID is the contract id for the Stellar asset.
	AssetContractID [32]byte
	// IsWrappedAsset is true if the contract id of the asset
	// is administered by the bridge contract
	IsWrappedAsset bool
	// Amount is the amount which will be transferred to the recipient.
	Amount int64
}

type StellarWithdrawalHandler struct {
	StellarBuilder         *stellar.Builder
	StellarSigner          *stellar.Signer
	StellarObserver        stellar.Observer
	WithdrawalWindow       time.Duration
	Converter              AssetConverter
	EthereumObserver       ethereum.Observer
	EthereumFinalityBuffer uint64
}

func (c *StellarWithdrawalHandler) CanWithdraw(deposit ethereum.Deposit) (StellarWithdrawalDetails, error) {
	assetContractID, isWrappedAsset, stellarAmount, err := c.Converter.ToStellar(deposit.Token.String(), deposit.Amount.String())
	if err != nil {
		return StellarWithdrawalDetails{}, err
	}

	destinationAccountID, err := strkey.Encode(
		strkey.VersionByteAccountID,
		deposit.Destination.Bytes(),
	)
	if err != nil {
		return StellarWithdrawalDetails{}, InvalidStellarRecipient
	}

	return StellarWithdrawalDetails{
		Deadline:        deposit.Time.Add(c.WithdrawalWindow),
		Recipient:       destinationAccountID,
		AssetContractID: assetContractID,
		IsWrappedAsset:  isWrappedAsset,
		Amount:          stellarAmount,
	}, nil
}

func (c *StellarWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	deposit, err := getEthereumDeposit(c.EthereumObserver, c.EthereumFinalityBuffer, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	details, err := c.CanWithdraw(deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	depositIDBytes, err := hex.DecodeString(deposit.ID)
	if err != nil {
		problem.Render(r.Context(), w, errors.Wrap(err, "error decoding deposit id"))
		return

	}
	tx, err := c.StellarBuilder.BuildTransaction(
		details.AssetContractID,
		details.IsWrappedAsset,
		sourceAccount,
		details.Recipient,
		amount.StringFromInt64(details.Amount),
		sequence,
		details.Deadline.Unix(),
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
