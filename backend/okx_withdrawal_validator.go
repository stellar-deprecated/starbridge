package backend

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stellar/go/support/render/problem"
	"github.com/stellar/starbridge/ethereum"
	"github.com/stellar/starbridge/store"
)

var (
	InvalidOkxRecipient = problem.P{
		Type:   "invalid_okx_recipient",
		Title:  "Invalid Okx Recipient",
		Status: http.StatusBadRequest,
		Detail: "The recipient of the deposit is not a valid Okx address.",
	}
)

// OkxWithdrawalValidator checks if it is possible to
// withdraw a deposit to the Stellar bridge account.
type OkxWithdrawalValidator struct {
	Observer          ethereum.Observer
	OkxFinalityBuffer uint64
	WithdrawalWindow  time.Duration
	Converter         AssetConverter
}

// OkxWithdrawalDetails includes metadata about the
// validation result.
type OkxWithdrawalDetails struct {
	// Deadline is the deadline for executing the withdrawal
	// transaction on Okx.
	Deadline time.Time
	// Recipient is the Okx address which should receive the
	// withdrawal.
	Recipient common.Address
	// Token is the address of the Okx tokens which will be
	// transferred to the recipient.
	Token common.Address
	// Amount is the amount of tokens which will be transferred to
	// the recipient.
	Amount *big.Int
}

func (s OkxWithdrawalValidator) CanWithdraw(ctx context.Context, deposit store.StellarDeposit) (OkxWithdrawalDetails, error) {
	if !common.IsHexAddress(deposit.Destination) {
		return OkxWithdrawalDetails{}, InvalidOkxRecipient
	}

	tokenAddress, amount, err := s.Converter.ToOkx(deposit.Asset, deposit.Amount)
	if err != nil {
		return OkxWithdrawalDetails{}, err
	}

	latest, err := s.Observer.GetLatestBlock(ctx)
	if err != nil {
		return OkxWithdrawalDetails{}, err
	}
	//if latest.Number <= s.EthereumFinalityBuffer {
	//	return EthereumWithdrawalDetails{}, EthereumNodeBehind
	//}

	latestFinalBlock, err := s.Observer.GetBlockByNumber(ctx, latest.Number-s.OkxFinalityBuffer)
	if err != nil {
		return OkxWithdrawalDetails{}, err
	}

	withdrawalDeadline := time.Unix(deposit.LedgerTime, 0).Add(s.WithdrawalWindow)
	if latestFinalBlock.Time.After(withdrawalDeadline) {
		return OkxWithdrawalDetails{}, WithdrawalWindowExpired
	}

	return OkxWithdrawalDetails{
		Deadline:  withdrawalDeadline,
		Recipient: common.HexToAddress(deposit.Destination),
		Token:     tokenAddress,
		Amount:    amount,
	}, nil
}
