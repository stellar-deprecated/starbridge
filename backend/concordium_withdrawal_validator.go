package backend

import (
	"context"
	"github.com/stellar/go/support/db"
	"github.com/stellar/starbridge/concordium"
	"math/big"
	"time"

	"github.com/stellar/starbridge/store"
)

//var (
//	InvalidEthereumRecipient = problem.P{
//		Type:   "invalid_ethereum_recipient",
//		Title:  "Invalid Ethereum Recipient",
//		Status: http.StatusBadRequest,
//		Detail: "The recipient of the deposit is not a valid Ethereum address.",
//	}
//	EthereumNodeBehind = problem.P{
//		Type:   "ethereum_node_behind",
//		Title:  "Ethereum Node Behind",
//		Status: http.StatusUnprocessableEntity,
//		Detail: "The ethereum node used by the validator is still catching up.",
//	}
//)

type ConcordiumWithdrawalValidator struct {
	Session          db.SessionInterface
	Observer         concordium.Observer
	WithdrawalWindow time.Duration
	Converter        AssetConverter
}

type ConcordiumWithdrawalDetails struct {
	Deadline  time.Time
	Recipient string
	Token     string
	Amount    *big.Int
	blockHash string
}

func (s ConcordiumWithdrawalValidator) CanWithdraw(ctx context.Context, deposit store.StellarDeposit) (ConcordiumWithdrawalDetails, error) {
	//if !common.IsHexAddress(deposit.Destination) {
	//	return EthereumWithdrawalDetails{}, InvalidEthereumRecipient
	//}

	tokenAddress, amount, err := s.Converter.ToConcordium(deposit.Asset, deposit.Amount)
	if err != nil {
		return ConcordiumWithdrawalDetails{}, err
	}

	//depositOn, err := s.Observer.GetDeposit(ctx, deposit.ID)
	//if err != nil {
	//	return ConcordiumWithdrawalDetails{}, err
	//}

	//latest, err := s.Observer.GetLatestBlock(ctx)
	//if err != nil {
	//	return ConcordiumWithdrawalDetails{}, err
	//}
	//if latest.Number <= s.EthereumFinalityBuffer {
	//	return EthereumWithdrawalDetails{}, EthereumNodeBehind
	//}

	//latestFinalBlock, err := s.Observer.GetBlockByNumber(ctx, latest.Number-s.EthereumFinalityBuffer)
	//if err != nil {
	//	return ConcordiumWithdrawalDetails{}, err
	//}

	withdrawalDeadline := time.Unix(deposit.LedgerTime, 0).Add(s.WithdrawalWindow)
	//if latestFinalBlock.Time.After(withdrawalDeadline) {
	//	return ConcordiumWithdrawalDetails{}, WithdrawalWindowExpired
	//}

	return ConcordiumWithdrawalDetails{
		Deadline:  withdrawalDeadline,
		Recipient: deposit.Destination,
		Token:     tokenAddress,
		Amount:    amount,
		//blockHash: depositOn.BlockHash,
	}, nil
}
