package transform

import (
	"fmt"
	"math"

	"github.com/stellar/starbridge/model"
)

// TODO this destination needs to be input from the input transaction on the remote chain from the memo or similar
var destinationAccountForStellar = "GCBAA5476KARHPDSU6WFQTPXQOWX3QMXU4LF7JVZ2ZMWJ4OQEL7ZMV6G" // var destinationSecretKey = "SALR2RNJG55BBWTML2MKO5CXG5QDI4ZTSVDIA53XDWOU7QPOAEQNYUE2"

// MapTxToChain converts a given transaction to the destination chain along with mapping the assets appropriately to the right contracts
func MapTxToChain(tx *model.Transaction, destinationChain *model.Chain) (*model.Transaction, error) {
	if destinationChain != model.ChainStellar {
		panic("can only convert to Stellar chain for now")
	}

	mappedAssetInfo, ok := destinationChain.AddressMappings[tx.AssetInfo]
	if !ok {
		return nil, fmt.Errorf("entry for input asset ('%s') did not exist, could not convert to mappedAssetInfo on destination chain", tx.AssetInfo.String())
	}

	// TODO set fee values here
	return &model.Transaction{
		Chain:                destinationChain,
		Hash:                 "",                           // TODO fill in converted tx hash
		Block:                0,                            // we don't have a block yet
		SeqNum:               destinationChain.NextNonce(), // TODO this can be asset specific too, but keeping it with a single nonce for now
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
