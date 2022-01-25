package integrations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stellar/starbridge/cmd/starbridge/model"
	"github.com/stellar/starbridge/contracts/gen/SimpleEscrowEvents"
)

// FetchEthTxByHash returns a model.Transaction
func FetchEthTxByHash(txHash string) (*model.Transaction, error) {
	infuraURL := "https://ropsten.infura.io/v3/a42ad89bbcec4ddca5c9abb60eb4a300"
	conn, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to infura (%s): %s", infuraURL, err)
	}
	defer conn.Close()

	ctx := context.Background()
	ethTxReceipt, err := conn.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("could not fetch transaction receipt '%s' by hash from network: %s", txHash, err)
	}
	ethTx, isPending, err := conn.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("could not fetch transaction '%s' by hash from network: %s", txHash, err)
	}

	tx, err := Ethereum2Transaction(conn, ethTxReceipt, ethTx, isPending)
	if err != nil {
		return nil, fmt.Errorf("could not convert txToString: %s", err)
	}
	return tx, nil
}

// Ethereum2Transaction makes a model.Transaction from an Ethereum Transaction
func Ethereum2Transaction(conn *ethclient.Client, txReceipt *types.Receipt, tx *types.Transaction, isPending bool) (*model.Transaction, error) {
	// we only allow funds sent to our contract for any asset
	if !IsMyContractAddress(tx.To().Hex()) {
		return nil, fmt.Errorf("unsupported receiver address '%s'", tx.To().Hex())
	}
	fromAddress, err := getFromAddress(tx)
	if err != nil {
		return nil, fmt.Errorf("unable to get From address: %s", err)
	}

	// pulled logic to read events from here: https://goethereumbook.org/event-read/ with many modifications to suit our smart contract
	// filter log events by a query
	query := ethereum.FilterQuery{
		FromBlock: txReceipt.BlockNumber,
		ToBlock:   txReceipt.BlockNumber,
		Addresses: []common.Address{
			common.HexToAddress(MY_ETHEREUM_CONTRACT_ADDRESS),
		},
	}
	logs, err := conn.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("unable to filter logs for contract address on block number %d: %w", txReceipt.BlockNumber.Int64(), err)
	}

	// parse filtered log events
	myAbi, err := abi.JSON(strings.NewReader(string(SimpleEscrowEvents.SimpleEscrowEventsABI)))
	if err != nil {
		return nil, fmt.Errorf("unable to read ABI: %s", err)
	}
	if len(logs) == 0 {
		return nil, fmt.Errorf("no event emitted for this tx hash")
	} else if len(logs) > 1 {
		return nil, fmt.Errorf("more than one event emitted for this tx hash")
	}
	vLog := logs[0]
	if len(vLog.Topics) != 1 {
		return nil, fmt.Errorf("we expect 1 topic entries for each event, for the event signature only")
	}
	event := PaymentEvent{}
	err = myAbi.UnpackIntoInterface(&event, eventName, vLog.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to unpack event into event type Payment: %s", err)
	}
	log.Printf("DEBUG - found event at txhash (%s): event.DestinationStellarAddress='%s', TokenAmount='%d', TokenContractAddress='%s'\n", vLog.TxHash.Hex(), event.DestinationStellarAddress, event.TokenAmount.Int64(), event.TokenContractAddress)

	// select asset using data from tx event
	assetInfo, ok := model.ChainEthereum.AllAssetMap[event.TokenContractAddress]
	if !ok {
		return nil, fmt.Errorf("found event with an unsupported contractAddress: %s", event.TokenContractAddress)
	}
	log.Printf("DEBUG - converted tokenContractAddress '%s' to assetInfo '%s'\n", event.TokenContractAddress, assetInfo)

	return &model.Transaction{
		Chain:     model.ChainEthereum,
		Hash:      txReceipt.TxHash.Hex(),
		Block:     txReceipt.BlockNumber.Uint64(),
		SeqNum:    tx.Nonce(),
		IsPending: isPending,
		From:      fromAddress,
		To:        tx.To().Hex(),
		AssetInfo: payableAsset,
		Amount:    tx.Value().Uint64(),
		Data: model.ContractData{
			EventName:                             eventName,
			TargetDestinationChain:                model.ChainStellar, // we have hard-coded this to Stellar for the MVP
			TargetDestinationAddressOnRemoteChain: event.DestinationStellarAddress,
			AssetInfo:                             assetInfo,
			Amount:                                uint64(event.TokenAmount.Int64()),
		},
		OriginalTx:           tx,
		AdditionalOriginalTx: []interface{}{txReceipt},
	}, nil
}

// getFromAddress gets the From address for an Ethereum transaction
func getFromAddress(tx *types.Transaction) (string, error) {
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if err != nil {
		return "", fmt.Errorf("could not get tx as message: %s", err)
	}
	return msg.From().Hex(), nil
}
