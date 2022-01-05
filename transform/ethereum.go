package transform

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stellar/starbridge/model"
)

// FetchEthTxByHash returns a model.Transaction
func FetchEthTxByHash(txHash string) (*model.Transaction, error) {
	conn, e := ethclient.Dial("https://ropsten.infura.io/v3/a42ad89bbcec4ddca5c9abb60eb4a300")
	if e != nil {
		return nil, fmt.Errorf("could not connect to infura testnet: %s", e)
	}

	ctx := context.Background()
	ethTxReceipt, e := conn.TransactionReceipt(ctx, common.HexToHash(txHash))
	if e != nil {
		return nil, fmt.Errorf("could not fetch transaction receipt '%s' by hash from network: %s", txHash, e)
	}
	ethTx, isPending, e := conn.TransactionByHash(ctx, common.HexToHash(txHash))
	if e != nil {
		return nil, fmt.Errorf("could not fetch transaction '%s' by hash from network: %s", txHash, e)
	}

	tx, e := Ethereum2Transaction(ethTxReceipt, ethTx, isPending)
	if e != nil {
		return nil, fmt.Errorf("could not convert txToString: %s", e)
	}
	return tx, nil
}

// Ethereum2Transaction makes a model.Transaction from an Ethereum Transaction
func Ethereum2Transaction(txReceipt *types.Receipt, tx *types.Transaction, isPending bool) (*model.Transaction, error) {
	fromAddress, e := getFromAddress(tx)
	if e != nil {
		return nil, fmt.Errorf("unable to get From address: %s", e)
	}

	var assetInfo *model.AssetInfo
	if txReceipt.ContractAddress.Hex() == model.AssetETH.ContractAddress {
		assetInfo = model.AssetETH
	} else {
		return nil, fmt.Errorf("unsupported contract address '%s' on Ethereum", txReceipt.ContractAddress.Hex())
	}

	return &model.Transaction{
		Chain:                model.ChainEthereum,
		Hash:                 txReceipt.TxHash.Hex(),
		Block:                txReceipt.BlockNumber.Uint64(),
		SeqNum:               tx.Nonce(),
		IsPending:            isPending,
		From:                 fromAddress,
		To:                   tx.To().Hex(),
		AssetInfo:            assetInfo,
		Amount:               tx.Value().Uint64(),
		OriginalTx:           tx,
		AdditionalOriginalTx: []interface{}{txReceipt},
	}, nil

	// sb.WriteString("Chain ID: " + tx.ChainId().String())
	// sb.WriteString("\nHash: " + txReceipt.TxHash.String())
	// sb.WriteString("\nIs Pending: " + fmt.Sprintf("%v", isPending))
	// sb.WriteString("\nBlock: " + txReceipt.BlockNumber.String())
	// sb.WriteString("\nNonce: " + fmt.Sprintf("%d", tx.Nonce()))
	// sb.WriteString("\nTx Index: " + fmt.Sprintf("%d", txReceipt.TransactionIndex))
	// sb.WriteString("\nContract Address: " + txReceipt.ContractAddress.String())
	// sb.WriteString("\nGas Price: " + tx.GasPrice().String())
	// sb.WriteString("\nCumulative Gas Used: " + fmt.Sprintf("%d", txReceipt.CumulativeGasUsed))
	// sb.WriteString("\nGas Used: " + fmt.Sprintf("%d", txReceipt.GasUsed))
	// sb.WriteString("\nStorage Size: " + txReceipt.Size().String())
	// sb.WriteString("\nCost: " + tx.Cost().String())
	// sb.WriteString("\nFrom: " + fromAddress)
	// sb.WriteString("\nTo: " + tx.To().Hex())
	// sb.WriteString("\nValue: " + tx.Value().String())
	// sb.WriteString("\nData: " + string(tx.Data()))
}

// getFromAddress gets the From address for an Ethereum transaction
func getFromAddress(tx *types.Transaction) (string, error) {
	msg, e := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(0))
	if e != nil {
		return "", fmt.Errorf("could not get tx as message: %s", e)
	}
	return msg.From().Hex(), nil
}
