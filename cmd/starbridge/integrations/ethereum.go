package integrations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stellar/starbridge/cmd/starbridge/model"
	"github.com/stellar/starbridge/contracts/gen/Counter"
)

// FetchEthTxByHash returns a model.Transaction
func FetchEthTxByHash(txHash string) (*model.Transaction, error) {
	infuraURL := "https://ropsten.infura.io/v3/a42ad89bbcec4ddca5c9abb60eb4a300"
	conn, e := ethclient.Dial(infuraURL)
	if e != nil {
		return nil, fmt.Errorf("could not connect to infura (%s): %s", infuraURL, e)
	}
	defer conn.Close()

	ctx := context.Background()
	ethTxReceipt, e := conn.TransactionReceipt(ctx, common.HexToHash(txHash))
	if e != nil {
		return nil, fmt.Errorf("could not fetch transaction receipt '%s' by hash from network: %s", txHash, e)
	}
	ethTx, isPending, e := conn.TransactionByHash(ctx, common.HexToHash(txHash))
	if e != nil {
		return nil, fmt.Errorf("could not fetch transaction '%s' by hash from network: %s", txHash, e)
	}

	tx, e := Ethereum2Transaction(conn, ethTxReceipt, ethTx, isPending)
	if e != nil {
		return nil, fmt.Errorf("could not convert txToString: %s", e)
	}
	return tx, nil
}

// Ethereum2Transaction makes a model.Transaction from an Ethereum Transaction
func Ethereum2Transaction(conn *ethclient.Client, txReceipt *types.Receipt, tx *types.Transaction, isPending bool) (*model.Transaction, error) {
	// we only allow funds sent to our contract for any asset
	if !IsMyContractAddress(tx.To().Hex()) {
		return nil, fmt.Errorf("unsupported receiver address '%s'", tx.To().Hex())
	}

	fromAddress, e := getFromAddress(tx)
	if e != nil {
		return nil, fmt.Errorf("unable to get From address: %s", e)
	}

	counter, e := Counter.NewCounter(common.HexToAddress(MY_ETHEREUM_CONTRACT_ADDRESS), conn)
	if e != nil {
		return nil, fmt.Errorf("unable to construct counter: %s", e)
	}
	cValue, e := counter.GetCount(nil)
	if e != nil {
		return nil, fmt.Errorf("unable to invoke GetCall(): %s", e)
	}
	log.Printf("value of counter is %d", cValue.Int64())

	myAbi, e := abi.JSON(strings.NewReader(Counter.CounterABI))
	if e != nil {
		return nil, fmt.Errorf("unable to read ABI: %s", e)
	}
	// TODO trying to get SHA3, see new code added in my_ethereum_contract to map fuinction names to sha hash values (4 bytes) so we can select correctly using fn sigs
	methodParams, e := myAbi.Unpack("incrementCounter", tx.Data()[4:])
	if e != nil {
		return nil, fmt.Errorf("unable to unpack ABI to get method params: %s", e)
	}
	log.Printf("method params sent in: %v", methodParams)

	// TODO convert the contract address to the correct token info - in this example it should give us AssetEthereum_USDC
	// ensure the contract exists in the list that we support
	// assetInfo, ok := model.ChainEthereum.AddressMappings[txReceipt.ContractAddress.Hex()]
	// if !ok {
	// 	return nil, fmt.Errorf("unsupported contract address '%s' on Ethereum", txReceipt.ContractAddress.Hex())
	// }
	// set this manually for now
	assetInfo := model.AssetEthereum_ETH

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
	msg, e := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if e != nil {
		return "", fmt.Errorf("could not get tx as message: %s", e)
	}
	return msg.From().Hex(), nil
}
