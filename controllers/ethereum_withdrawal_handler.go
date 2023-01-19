package controllers

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stellar/go/support/render/problem"

	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/stellar"
)

var (
	InvalidEthereumRecipient = problem.P{
		Type:   "invalid_ethereum_recipient",
		Title:  "Invalid Ethereum Recipient",
		Status: http.StatusBadRequest,
		Detail: "The recipient of the deposit is not a valid Ethereum address.",
	}
	EthereumNodeBehind = problem.P{
		Type:   "ethereum_node_behind",
		Title:  "Ethereum Node Behind",
		Status: http.StatusUnprocessableEntity,
		Detail: "The ethereum node used by the validator is still catching up.",
	}
	WithdrawalWindowExpired = problem.P{
		Type:   "withdrawal_window_expired",
		Title:  "Withdrawal Window Expired",
		Status: http.StatusBadRequest,
		Detail: "The withdrawal window has expired. Only refunds are allowed at this point.",
	}
)

// EthereumWithdrawalDetails includes metadata about the
// validation result.
type EthereumWithdrawalDetails struct {
	// Deadline is the deadline for executing the withdrawal
	// transaction on Ethereum.
	Deadline time.Time
	// Recipient is the Ethereum address which should receive the
	// withdrawal.
	Recipient common.Address
	// Token is the address of the Ethereum tokens which will be
	// transferred to the recipient.
	Token common.Address
	// Amount is the amount of tokens which will be transferred to
	// the recipient.
	Amount *big.Int
}

type EthereumWithdrawalHandler struct {
	EthereumObserver       ethereum.Observer
	StellarObserver        stellar.Observer
	EthereumSigner         ethereum.Signer
	EthereumFinalityBuffer uint64
	WithdrawalWindow       time.Duration
	Converter              AssetConverter
}

func (c *EthereumWithdrawalHandler) CanWithdraw(deposit stellar.Deposit) (EthereumWithdrawalDetails, error) {
	if !common.IsHexAddress(deposit.Destination) {
		return EthereumWithdrawalDetails{}, InvalidEthereumRecipient
	}

	tokenAddress, amount, err := c.Converter.ToEthereum(deposit.Token, deposit.Amount)
	if err != nil {
		return EthereumWithdrawalDetails{}, err
	}

	return EthereumWithdrawalDetails{
		Deadline:  deposit.Time.Add(c.WithdrawalWindow),
		Recipient: common.HexToAddress(deposit.Destination),
		Token:     tokenAddress,
		Amount:    amount,
	}, nil
}

func (c *EthereumWithdrawalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	deposit, err := getStellarDeposit(c.StellarObserver, r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	details, err := c.CanWithdraw(deposit)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	sig, err := c.EthereumSigner.SignWithdrawal(
		common.HexToHash(deposit.ID),
		details.Deadline.Unix(),
		details.Recipient,
		details.Token,
		details.Amount,
	)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	responseBytes, err := json.Marshal(EthereumSignatureResponse{
		Address:    c.EthereumSigner.Address().String(),
		Signature:  hex.EncodeToString(sig),
		DepositID:  deposit.ID,
		Expiration: details.Deadline.Unix(),
		Token:      details.Token.String(),
		Amount:     details.Amount.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, _ = w.Write(responseBytes)
	}
}
