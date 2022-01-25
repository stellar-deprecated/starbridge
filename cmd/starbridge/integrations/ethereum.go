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

	// pulled logic to read events from here: https://goethereumbook.org/event-read/ with many modifications to suit our smart contract
	// filter log events by a query
	query := ethereum.FilterQuery{
		FromBlock: txReceipt.BlockNumber,
		ToBlock:   txReceipt.BlockNumber,
		Addresses: []common.Address{
			common.HexToAddress(MY_ETHEREUM_CONTRACT_ADDRESS),
		},
	}
	logs, e := conn.FilterLogs(context.Background(), query)
	if e != nil {
		return nil, fmt.Errorf("unable to filter logs for contract address on block number %d: %w", txReceipt.BlockNumber.Int64(), e)
	}

	// parse filtered log events
	var assetInfo *model.AssetInfo
	var eventTokenAmount uint64
	myAbi, e := abi.JSON(strings.NewReader(string(SimpleEscrowEvents.SimpleEscrowEventsABI)))
	if e != nil {
		return nil, fmt.Errorf("unable to read ABI: %s", e)
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
	e = myAbi.UnpackIntoInterface(&event, eventName, vLog.Data)
	if e != nil {
		return nil, fmt.Errorf("unable to unpack event into event type Payment: %s", e)
	}
	// set values from log event
	// TODO make this use the map in the Chain directly, need to store it by keccak256 hash value too
	//     assetInfo, ok := model.ChainEthereum.AddressMappings[txReceipt.ContractAddress.Hex()]
	eventDestinationStellarAddress := event.DestinationStellarAddress
	eventTokenAmount = uint64(event.TokenAmount.Int64())
	if event.TokenContractAddress == model.AssetEthereum_ETH.ContractAddress {
		log.Printf("DEBUG - found event with ethContractAddress at txhash (%s): event.DestinationStellarAddress='%s', TokenAmount='%d'\n", vLog.TxHash.Hex(), event.DestinationStellarAddress, event.TokenAmount.Int64())
		assetInfo = model.AssetEthereum_ETH
	} else if event.TokenContractAddress == model.AssetEthereum_USDC.ContractAddress {
		log.Printf("DEBUG - found event with usdcContractAddress at txhash (%s): event.DestinationStellarAddress='%s', TokenAmount='%d'\n", vLog.TxHash.Hex(), event.DestinationStellarAddress, event.TokenAmount.Int64())
		assetInfo = model.AssetEthereum_USDC
	} else {
		return nil, fmt.Errorf("found event with an unsupported contractAddress: %s", event.TokenContractAddress)
	}
	contractData := model.ContractData{
		EventName:                             eventName,
		TargetDestinationChain:                model.ChainStellar, // we have hard-coded this to Stellar for the MVP
		TargetDestinationAddressOnRemoteChain: eventDestinationStellarAddress,
		AssetInfo:                             assetInfo,
		Amount:                                eventTokenAmount,
	}

	// TODO parse Stellar destination address and see if it is valid

	return &model.Transaction{
		Chain:                model.ChainEthereum,
		Hash:                 txReceipt.TxHash.Hex(),
		Block:                txReceipt.BlockNumber.Uint64(),
		SeqNum:               tx.Nonce(),
		IsPending:            isPending,
		From:                 fromAddress,
		To:                   tx.To().Hex(),
		AssetInfo:            payableAsset,
		Amount:               tx.Value().Uint64(),
		Data:                 contractData,
		OriginalTx:           tx,
		AdditionalOriginalTx: []interface{}{txReceipt},
	}, nil
}

// getFromAddress gets the From address for an Ethereum transaction
func getFromAddress(tx *types.Transaction) (string, error) {
	msg, e := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if e != nil {
		return "", fmt.Errorf("could not get tx as message: %s", e)
	}
	return msg.From().Hex(), nil
}
