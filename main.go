package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	txHash := "0x13070f64d40f22cd10c5bf9972767b67406ed3d818a50f82b1409289dcaa1aec"
	txReceiptJsonString := "{\"root\":\"0x\",\"status\":\"0x1\",\"cumulativeGasUsed\":\"0x2380ba\",\"logsBloom\":\"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\",\"logs\":[],\"transactionHash\":\"0x13070f64d40f22cd10c5bf9972767b67406ed3d818a50f82b1409289dcaa1aec\",\"contractAddress\":\"0x0000000000000000000000000000000000000000\",\"gasUsed\":\"0x5208\",\"blockHash\":\"0xf8d106bdd68f4d7cceaa6458c37cc03b0675d672ab63a68d33f4574c27c24e7f\",\"blockNumber\":\"0xb35b44\",\"transactionIndex\":\"0x21\"}"
	txJsonString := ""
	readFromNetwork := true

	var e error
	var txReceipt *types.Receipt
	var tx *types.Transaction
	var isPending bool
	if readFromNetwork {
		conn, e := ethclient.Dial("https://ropsten.infura.io/v3/a42ad89bbcec4ddca5c9abb60eb4a300")
		if e != nil {
			log.Fatal(fmt.Errorf("could not connect to infura testnet: %s", e))
		}

		ctx := context.Background()
		txReceipt, e = conn.TransactionReceipt(ctx, common.HexToHash(txHash))
		if e != nil {
			log.Fatal(fmt.Errorf("could not fetch transaction receipt '%s' by hash from network: %s", txHash, e))
		}
		tx, isPending, e = conn.TransactionByHash(ctx, common.HexToHash(txHash))
		if e != nil {
			log.Fatal(fmt.Errorf("could not fetch transaction '%s' by hash from network: %s", txHash, e))
		}
	} else {
		e = json.Unmarshal([]byte(txReceiptJsonString), &txReceipt)
		if e != nil {
			log.Fatal(fmt.Errorf("could not fetch transaction receipt '%s' by hash from saved json: %s", txHash, e))
		}

		e = json.Unmarshal([]byte(txJsonString), &tx)
		if e != nil {
			log.Fatal(fmt.Errorf("could not fetch transaction '%s' by hash from saved json: %s", txHash, e))
		}
	}

	txString, e := txToString(txReceipt, tx, isPending)
	if e != nil {
		log.Fatal(fmt.Errorf("could not convert txToString: %s", e))
	}
	fmt.Println(txString)
}

func txToString(txReceipt *types.Receipt, tx *types.Transaction, isPending bool) (string, error) {
	fromAddress, e := getFromAddress(tx)
	if e != nil {
		return "", fmt.Errorf("unable to get From address: %s", e)
	}

	sb := strings.Builder{}
	sb.WriteString("Chain ID: " + tx.ChainId().String())
	sb.WriteString("\nHash: " + txReceipt.TxHash.String())
	sb.WriteString("\nIs Pending: " + fmt.Sprintf("%v", isPending))
	sb.WriteString("\nBlock: " + txReceipt.BlockNumber.String())
	sb.WriteString("\nNonce: " + fmt.Sprintf("%d", tx.Nonce()))
	sb.WriteString("\nTx Index: " + fmt.Sprintf("%d", txReceipt.TransactionIndex))
	sb.WriteString("\nContract Address: " + txReceipt.ContractAddress.String())
	sb.WriteString("\nGas Price: " + tx.GasPrice().String())
	sb.WriteString("\nCumulative Gas Used: " + fmt.Sprintf("%d", txReceipt.CumulativeGasUsed))
	sb.WriteString("\nGas Used: " + fmt.Sprintf("%d", txReceipt.GasUsed))
	sb.WriteString("\nStorage Size: " + txReceipt.Size().String())
	sb.WriteString("\nCost: " + tx.Cost().String())
	sb.WriteString("\nFrom: " + fromAddress)
	sb.WriteString("\nTo: " + tx.To().Hex())
	sb.WriteString("\nValue: " + tx.Value().String())
	sb.WriteString("\nData: " + string(tx.Data()))
	return sb.String(), nil
}

func getFromAddress(tx *types.Transaction) (string, error) {
	msg, e := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(0))
	if e != nil {
		return "", fmt.Errorf("could not get tx as message: %s", e)
	}
	return msg.From().Hex(), nil
}
