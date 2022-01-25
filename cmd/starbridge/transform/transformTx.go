package transform

import (
	"fmt"
	"math"

	"github.com/stellar/starbridge/cmd/starbridge/integrations"
	"github.com/stellar/starbridge/cmd/starbridge/model"
)

// TODO this destination needs to be input from the input transaction on the remote chain from the memo or similar
var destinationAccountForStellar = "GBUKJ5TXOBDB5SKVCEV27OYLHUL2ZD3OITMVFNT7MZ4LPO3EECTLADCU" // var destinationSecretKey = "SB6HKIGE6KKAOMSFY7G7VFIKNKMSQ7QRASFI65CVJLTZ7SUPLY4FAZ3L"

// MapTxToChain converts a given transaction to the destination chain along with mapping the assets appropriately to the right contracts
func MapTxToChain(tx *model.Transaction, destinationChain *model.Chain) (*model.Transaction, error) {
	if destinationChain != model.ChainStellar {
		return nil, fmt.Errorf("can only convert to Stellar chain for now")
	}

	mappedAssetInfo, ok := destinationChain.AddressMappings[tx.AssetInfo.MapKey()]
	if !ok {
		return nil, fmt.Errorf("entry for input asset ('%s') did not exist, could not convert to mappedAssetInfo on destination chain with addressMappings: %+v", tx.AssetInfo.String(), destinationChain.AddressMappings)
	}

	nextNonce, e := destinationChain.NextNonce(integrations.GetSourceAccount())
	if e != nil {
		return nil, fmt.Errorf("cannot get next nonce: %s", e)
	}

	// TODO set fee values here
	return &model.Transaction{
		Chain:                destinationChain,
		Hash:                 "", // TODO fill in converted tx hash
		Block:                0,  // we don't have a block yet
		SeqNum:               nextNonce,
		IsPending:            true,
		From:                 mappedAssetInfo.ContractAddress, // we want to send from the Stellar escrow account for now
		To:                   destinationAccountForStellar,    // TODO fetch from the input transaction
		AssetInfo:            mappedAssetInfo,
		Amount:               amountUsingDecimals(tx.AssetInfo.Decimals, mappedAssetInfo.Decimals, tx.Amount),
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
