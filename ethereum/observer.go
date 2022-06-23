package ethereum

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stellar/starbridge/solidity-go"
)

var (
	// ErrLogNotFound is returned by GetDeposit when the log
	// with the given block index cannot be found
	ErrLogNotFound = fmt.Errorf("log not found")
	// ErrLogNotFromBridge is returned by GetDeposit when the log
	// with the given block index is not emitted from the bridge contract
	ErrLogNotFromBridge = fmt.Errorf("log is not from bridge")
)

// Block represents an ethereum block
type Block struct {
	// Number is the sequence number of the block
	Number uint64
	// Time is the timestamp when the block was executed
	Time time.Time
}

// RequestStatus is the status of a withdrawal on the
// bridge contract
type RequestStatus struct {
	// Fulfilled is true if the withdrawal was executed
	Fulfilled bool
	// BlockNumber is the latst block at the time the
	// request status was queried
	BlockNumber uint64
}

// Deposit is a deposit to the bridge smart contract
type Deposit struct {
	// Token is the address (0x0 in the case that eth was deposited)
	// of the tokens which were deposited to the bridge
	Token common.Address
	// Sender is the address of the account which deposited the tokens
	Sender common.Address
	// Destination is the intended recipient of the bridge transfer
	Destination *big.Int
	// Amount is the amount of tokens which were deposited to the bridge
	// contract
	Amount *big.Int
	// TxHash is the hash of the transaction containing the deposit
	TxHash common.Hash
	// LogIndex is the log index within the ethereum block of the deposit event
	// emitted by the bridge contract
	LogIndex uint
	// BlockNumber is the sequence number of the block containing the deposit
	// transaction
	BlockNumber uint64
	// Time is the timestamp of the deposit
	Time time.Time
}

// ID returns a unique id for the deposit
func (d Deposit) ID() common.Hash {
	logIndexBytes := [32]byte{}
	binary.PutUvarint(logIndexBytes[:], uint64(d.LogIndex))
	return crypto.Keccak256Hash(d.TxHash[:], logIndexBytes[:])
}

// Observer is used to inspect the ethereum blockchain to
// for all information relevant to bridge interactions
type Observer struct {
	client        *ethclient.Client
	filterer      *solidity.BridgeFilterer
	caller        *solidity.BridgeCaller
	bridgeAddress common.Address
}

// NewObserver constructs a new Observer instance
func NewObserver(client *ethclient.Client, bridgeAddress common.Address) (Observer, error) {
	caller, err := solidity.NewBridgeCaller(bridgeAddress, client)
	if err != nil {
		return Observer{}, err
	}
	filterer, err := solidity.NewBridgeFilterer(bridgeAddress, client)
	if err != nil {
		return Observer{}, err
	}

	return Observer{
		client:        client,
		filterer:      filterer,
		caller:        caller,
		bridgeAddress: bridgeAddress,
	}, nil
}

// GetDeposit returns a Deposit instance identified by the given transaction
// hash and log index
func (o Observer) GetDeposit(
	ctx context.Context, txHash string, logIndex uint,
) (Deposit, error) {
	receipt, err := o.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return Deposit{}, err
	}
	var log *types.Log
	for _, l := range receipt.Logs {
		if l.Index == logIndex {
			log = l
			break
		}
	}
	if log == nil {
		return Deposit{}, ErrLogNotFound
	}
	if log.Address != o.bridgeAddress {
		return Deposit{}, ErrLogNotFromBridge
	}

	header, err := o.client.HeaderByHash(ctx, log.BlockHash)
	if err != nil {
		return Deposit{}, err
	}

	erc20Event, err := o.filterer.ParseDepositERC20(*log)
	if err == nil {
		return Deposit{
			Token:       erc20Event.Token,
			Sender:      erc20Event.Sender,
			Destination: erc20Event.Destination,
			Amount:      erc20Event.Amount,
			TxHash:      log.TxHash,
			LogIndex:    logIndex,
			BlockNumber: log.BlockNumber,
			Time:        time.Unix(int64(header.Time), 0),
		}, nil
	}

	ethEvent, err := o.filterer.ParseDepositETH(*log)
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{
		Token:       common.Address{},
		Sender:      ethEvent.Sender,
		Destination: ethEvent.Destination,
		Amount:      ethEvent.Amount,
		TxHash:      log.TxHash,
		LogIndex:    logIndex,
		BlockNumber: log.BlockNumber,
		Time:        time.Unix(int64(header.Time), 0),
	}, nil
}

// GetLatestBlock returns the latest ethereum block
func (o Observer) GetLatestBlock(ctx context.Context) (Block, error) {
	return o.getBlockByNumber(ctx, nil)
}

// GetBlockByNumber finds an ethereum block by its sequence number
func (o Observer) GetBlockByNumber(ctx context.Context, blockNumber uint64) (Block, error) {
	bn := &big.Int{}
	bn.SetUint64(blockNumber)
	return o.getBlockByNumber(ctx, bn)
}

func (o Observer) getBlockByNumber(ctx context.Context, number *big.Int) (Block, error) {
	header, err := o.client.HeaderByNumber(ctx, number)
	if err != nil {
		return Block{}, err
	}
	return Block{
		Number: header.Number.Uint64(),
		Time:   time.Unix(int64(header.Time), 0),
	}, nil
}

// GetRequestStatus calls the requestStatus() view function on the bridge contract
// to determine the status of a bridge withdrawal
func (o Observer) GetRequestStatus(ctx context.Context, requestID common.Hash) (RequestStatus, error) {
	fulfilled, blockNumber, err := o.caller.RequestStatus(&bind.CallOpts{Context: ctx}, requestID)
	if err != nil {
		return RequestStatus{}, err
	}
	return RequestStatus{
		Fulfilled:   fulfilled,
		BlockNumber: blockNumber.Uint64(),
	}, nil
}
