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
)

// EthereumWithdrawalValidator checks if it is possible to
// withdraw a deposit to the Stellar bridge account.
type EthereumWithdrawalValidator struct {
	Observer               ethereum.Observer
	EthereumFinalityBuffer uint64
	WithdrawalWindow       time.Duration
	Converter              AssetConverter
}

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

func (s EthereumWithdrawalValidator) CanWithdraw(ctx context.Context, deposit store.StellarDeposit) (EthereumWithdrawalDetails, error) {
	if !common.IsHexAddress(deposit.Destination) {
		return EthereumWithdrawalDetails{}, InvalidEthereumRecipient
	}

	tokenAddress, amount, err := s.Converter.ToEthereum(deposit.Asset, deposit.Amount)
	if err != nil {
		return EthereumWithdrawalDetails{}, err
	}

	latest, err := s.Observer.GetLatestBlock(ctx)
	if err != nil {
		return EthereumWithdrawalDetails{}, err
	}
	//if latest.Number <= s.EthereumFinalityBuffer {
	//	return EthereumWithdrawalDetails{}, EthereumNodeBehind
	//}

	latestFinalBlock, err := s.Observer.GetBlockByNumber(ctx, latest.Number-s.EthereumFinalityBuffer)
	if err != nil {
		return EthereumWithdrawalDetails{}, err
	}

	withdrawalDeadline := time.Unix(deposit.LedgerTime, 0).Add(s.WithdrawalWindow)
	if latestFinalBlock.Time.After(withdrawalDeadline) {
		return EthereumWithdrawalDetails{}, WithdrawalWindowExpired
	}

	return EthereumWithdrawalDetails{
		Deadline:  withdrawalDeadline,
		Recipient: common.HexToAddress(deposit.Destination),
		Token:     tokenAddress,
		Amount:    amount,
	}, nil
}
