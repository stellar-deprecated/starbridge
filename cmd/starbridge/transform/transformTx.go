package transform

import (
	"fmt"
	"math"

	"github.com/stellar/starbridge/cmd/starbridge/integrations"
	"github.com/stellar/starbridge/cmd/starbridge/model"
)

// MapTxToChain converts a given transaction to the destination chain along with mapping the assets appropriately to the right contracts
func MapTxToChain(tx *model.Transaction) (*model.Transaction, error) {
	if tx.Chain != model.ChainEthereum {
		return nil, fmt.Errorf("can only convert from Ethereum chain for now")
	}
	if tx.Data.TargetDestinationChain != model.ChainStellar {
		return nil, fmt.Errorf("can only convert to Stellar chain for now")
	}

	mappedAssetInfo, ok := tx.Data.TargetDestinationChain.AddressMappings[tx.Data.AssetInfo.MapKey()]
	if !ok {
		return nil, fmt.Errorf("entry for input asset ('%s') did not exist, could not convert to mappedAssetInfo on destination chain with addressMappings: %+v", tx.Data.AssetInfo.String(), tx.Data.TargetDestinationChain.AddressMappings)
	}

	// TODO the source account should maybe come directly from the chain but adding here to avoid an import cycle. Need to move files around
	nextNonce, e := tx.Data.TargetDestinationChain.NextNonce(integrations.GetSourceAccount())
	if e != nil {
		return nil, fmt.Errorf("cannot get next nonce: %s", e)
	}

	// TODO set fee values here
	return &model.Transaction{
		Chain:                tx.Data.TargetDestinationChain,
		Hash:                 "", // TODO fill in converted tx hash
		Block:                0,  // we don't have a block yet
		SeqNum:               nextNonce,
		IsPending:            true,
		From:                 mappedAssetInfo.ContractAddress,               // escrow account is always the account that will send the payment
		To:                   tx.Data.TargetDestinationAddressOnRemoteChain, // this is where we do the conversion of the contract data to the To account
		AssetInfo:            mappedAssetInfo,
		Amount:               amountUsingDecimals(tx.Data.AssetInfo.Decimals, mappedAssetInfo.Decimals, tx.Data.Amount), // this is where we do the conversion of the contract data to the Amount value
		Data:                 tx.Data,
		OriginalTx:           tx,
		AdditionalOriginalTx: nil,
	}, nil
}

// AmountUsingDecimals is a helper that converts decimal values for us
func amountUsingDecimals(fromDecimals int, newDecimals int, amount uint64) uint64 {
	var amountExponent int = newDecimals - fromDecimals
	if amountExponent > 0 {
		return amount * uint64(math.Pow10(amountExponent))
	} else if amountExponent < 0 {
		return uint64(amount / uint64(math.Pow10(-amountExponent)))
	}
	return amount
}
