package ethereum

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"

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
	// ErrLogNotDepositEvent is returned by GetDeposit when the log
	// with the given block index is not a DepositETH or DepositERC20 event
	ErrLogNotDepositEvent = fmt.Errorf("log is not a deposit event")
	// ErrTxHashNotFound is returned by GetDeposit when the given transaction
	// hash is not found
	ErrTxHashNotFound = fmt.Errorf("deposit tx hash not found")
)

// IsInvalidGetDepositRequest returns true if the given error
// from GetDeposit indicates that the provided transaction hash
// or log index is invalid
func IsInvalidGetDepositRequest(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrLogNotFound) ||
		errors.Is(err, ErrLogNotFromBridge) ||
		errors.Is(err, ErrLogNotDepositEvent)
}

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

// DepositID returns a globally unique id for a given deposit
func DepositID(txHash string, logIndex uint) string {
	hash := common.HexToHash(txHash)
	logIndexBytes := [32]byte{}
	binary.PutUvarint(logIndexBytes[:], uint64(logIndex))
	id := crypto.Keccak256Hash(hash[:], logIndexBytes[:])
	return hex.EncodeToString(id.Bytes())
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
func NewObserver(client *ethclient.Client, bridgeAddress string) (Observer, error) {
	if !common.IsHexAddress(bridgeAddress) {
		return Observer{}, fmt.Errorf("%v is not a valid ethereum address", bridgeAddress)
	}
	bridgeAddressParsed := common.HexToAddress(bridgeAddress)

	caller, err := solidity.NewBridgeCaller(bridgeAddressParsed, client)
	if err != nil {
		return Observer{}, err
	}
	filterer, err := solidity.NewBridgeFilterer(bridgeAddressParsed, client)
	if err != nil {
		return Observer{}, err
	}

	return Observer{
		client:        client,
		filterer:      filterer,
		caller:        caller,
		bridgeAddress: bridgeAddressParsed,
	}, nil
}

// GetDeposit returns a Deposit instance identified by the given transaction
// hash and log index
func (o Observer) GetDeposit(
	ctx context.Context, txHash string, logIndex uint,
) (Deposit, error) {
	receipt, err := o.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		if err == ethereum.NotFound {
			return Deposit{}, ErrTxHashNotFound
		}
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

	event, err := o.filterer.ParseDeposit(*log)
	if err != nil {
		return Deposit{}, ErrLogNotDepositEvent
	}
	return Deposit{
		Token:       event.Token,
		Sender:      event.Sender,
		Destination: event.Destination,
		Amount:      event.Amount,
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

// GetDomainSeparator calls the domainSeparator public attribute on the bridge contract
// and returns its value
func (o Observer) GetDomainSeparator(ctx context.Context) ([32]byte, error) {
	return o.caller.DomainSeparator(&bind.CallOpts{Context: ctx})
}
