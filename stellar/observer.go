package stellar

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/clients/stellarcore"
	stellarcoreproto "github.com/stellar/go/protocols/stellarcore"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrEventNotFound is returned by GetDeposit when the deposit event cannot be found
	ErrEventNotFound = fmt.Errorf("event not found")
	// ErrEventNotFromBridge is returned by GetDeposit when the event
	// is not emitted from the bridge contract
	ErrEventNotFromBridge = fmt.Errorf("event is not from bridge")
	// ErrNotDepositEvent is returned by GetDeposit when the event is not a
	// valid deposit event
	ErrNotDepositEvent = fmt.Errorf("log is not a deposit event")
	// ErrTxHashNotFound is returned by GetDeposit when the given transaction
	// hash is not found
	ErrTxHashNotFound = fmt.Errorf("deposit tx hash not found")
	// ErrSenderIsNotAccount is returned by GetDeposit when the sender of the deposit
	// is not an account
	ErrSenderIsNotAccount = fmt.Errorf("sender of deposit is not account")
)

// IsInvalidGetDepositRequest returns true if the given error
// from GetDeposit indicates that the provided transaction hash
// or log index is invalid
func IsInvalidGetDepositRequest(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrEventNotFound) ||
		errors.Is(err, ErrEventNotFromBridge) ||
		errors.Is(err, ErrNotDepositEvent) ||
		errors.Is(err, ErrSenderIsNotAccount)
}

// RequestStatus is the status of a withdrawal on the
// bridge contract
type RequestStatus struct {
	// Fulfilled is true if the withdrawal was executed
	Fulfilled bool
	// LedgerSequence is the latest ledger sequence at the time the
	// request status was queried
	LedgerSequence uint32
	// Time is the timestamp of the latest ledger at the time the
	// request status was queried
	CloseTime time.Time
}

// Deposit is a deposit to the bridge smart contract
type Deposit struct {
	// ID is the globally unique id for a given deposit
	ID string
	// Token is the contract id
	// of the tokens which were deposited to the bridge
	Token [32]byte
	// IsWrappedAsset is true if the contract id of the asset
	// is administered by the bridge contract
	IsWrappedAsset bool
	// Sender is the account which deposited the tokens
	Sender string
	// Destination is the intended recipient of the bridge transfer
	Destination string
	// Amount is the amount of tokens which were deposited to the bridge
	// contract
	Amount string
	// TxHash is the hash of the transaction containing the deposit
	TxHash string
	// OperationIndex is the index within the operations for the transaction
	// corresponding to the operation containing the deposit
	OperationIndex uint
	// EventIndex is the index within the operation events for the deposit event
	EventIndex uint
	// Time is the timestamp of the deposit
	Time time.Time
}

func depositID(txHash string, operationIndex, eventIndex uint) string {
	hash := common.HexToHash(txHash)
	operationIndexBytes := [32]byte{}
	binary.PutUvarint(operationIndexBytes[:], uint64(operationIndex))
	logIndexBytes := [32]byte{}
	binary.PutUvarint(logIndexBytes[:], uint64(eventIndex))
	id := crypto.Keccak256Hash(hash[:], operationIndexBytes[:], logIndexBytes[:])
	return hex.EncodeToString(id.Bytes())
}

// Observer is used to inspect the ethereum blockchain to
// for all information relevant to bridge interactions
type Observer struct {
	bridgeContractID [32]byte
	horizonClient    *horizonclient.Client
	coreClient       *stellarcore.Client
}

// NewObserver constructs a new Observer instance
func NewObserver(bridgeContractID [32]byte, horizonClient *horizonclient.Client, coreClient *stellarcore.Client) Observer {
	return Observer{
		bridgeContractID: bridgeContractID,
		horizonClient:    horizonClient,
		coreClient:       coreClient,
	}
}

// GetDeposit returns a Deposit instance identified by the given transaction
// hash and log index
func (o Observer) GetDeposit(
	ctx context.Context, txHash string, opIndex, eventIndex uint,
) (Deposit, error) {
	tx, err := o.horizonClient.TransactionDetail(txHash)
	if err != nil {
		if horizonclient.IsNotFoundError(err) {
			return Deposit{}, ErrTxHashNotFound
		}
		return Deposit{}, err
	}

	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXdr, &meta); err != nil {
		// Invalid meta back. Eek!
		return Deposit{}, err
	}

	v3, ok := meta.GetV3()
	if !ok {
		return Deposit{}, ErrEventNotFound
	}
	if opIndex >= uint(len(v3.Events)) || eventIndex >= uint(len(v3.Events[opIndex].Events)) {
		return Deposit{}, ErrEventNotFound
	}
	event := v3.Events[opIndex].Events[eventIndex]
	if event.ContractId == nil || *event.ContractId != o.bridgeContractID {
		return Deposit{}, ErrEventNotFromBridge
	}

	depositSym := xdr.ScSymbol("deposit")
	eventBody := event.Body.MustV0()
	if len(eventBody.Topics) != 4 ||
		eventBody.Topics[0].Equals(xdr.ScVal{Type: xdr.ScValTypeScvSymbol, Sym: &depositSym}) {
		return Deposit{}, ErrEventNotFromBridge
	}

	tokenBytes := eventBody.Topics[1].MustObj().MustBin()
	if len(tokenBytes) != 32 {
		return Deposit{}, ErrEventNotFromBridge
	}
	var token [32]byte
	copy(token[:], tokenBytes)

	identifierVec := eventBody.Topics[2].MustObj().MustVec()
	accountSym := xdr.ScSymbol("Account")
	if identifierVec[0].Equals(xdr.ScVal{Type: xdr.ScValTypeScvSymbol, Sym: &accountSym}) {
		return Deposit{}, ErrSenderIsNotAccount
	}

	sender := identifierVec[1].MustObj().MustAccountId()
	destination := eventBody.Topics[3].MustObj().MustAccountId()
	data := eventBody.Data.MustObj().MustVec()
	amounti128 := data[0].MustObj().MustI128()
	lo := xdr.Int64(amounti128.Lo)
	if amounti128.Hi > 0 || lo < 0 {
		return Deposit{}, fmt.Errorf("deposit amount is too high")
	}
	isWrappedAsset := data[1].MustIc() == xdr.ScStaticScsTrue

	return Deposit{
		ID:             depositID(txHash, opIndex, eventIndex),
		Token:          token,
		IsWrappedAsset: isWrappedAsset,
		Sender:         sender.Address(),
		Destination:    destination.Address(),
		Amount:         amount.String(lo),
		TxHash:         tx.Hash,
		OperationIndex: opIndex,
		EventIndex:     eventIndex,
		Time:           tx.LedgerCloseTime,
	}, nil
}

func functionNameParam(name string) xdr.ScVal {
	contractFnParameterSym := xdr.ScSymbol(name)
	return xdr.ScVal{
		Type: xdr.ScValTypeScvSymbol,
		Sym:  &contractFnParameterSym,
	}
}

func bytes32ContractParam(contractID xdr.Hash) xdr.ScVal {
	contractIdBytes := contractID[:]
	contractIdParameterObj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoBytes,
		Bin:  &contractIdBytes,
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &contractIdParameterObj,
	}
}

func preflight(coreClient *stellarcore.Client, invokeHostFn *txnbuild.InvokeHostFunction) (stellarcoreproto.PreflightResponse, error) {
	opXDR, err := invokeHostFn.BuildXDR()
	if err != nil {
		return stellarcoreproto.PreflightResponse{}, err
	}

	invokeHostFunctionOp := opXDR.Body.MustInvokeHostFunctionOp()

	response, err := coreClient.Preflight(
		context.Background(),
		invokeHostFn.SourceAccount,
		invokeHostFunctionOp,
	)
	if err != nil {
		return response, err
	}

	if response.Status != stellarcoreproto.PreflightStatusOk {
		return response, fmt.Errorf("status is not ok: %v", response.Detail)
	}
	return response, nil
}

// GetRequestStatus calls the status() function on the bridge contract
// to determine the status of a bridge withdrawal
func (o Observer) GetRequestStatus(ctx context.Context, requestID [32]byte) (RequestStatus, error) {
	invokeStatus := &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeArgs: &xdr.ScVec{
				bytes32ContractParam(o.bridgeContractID),
				functionNameParam("status"),
				bytes32ContractParam(requestID),
			},
		},
	}

	response, err := preflight(o.coreClient, invokeStatus)
	if err != nil {
		return RequestStatus{}, err
	}
	var status xdr.ScVal
	if err = xdr.SafeUnmarshalBase64(response.Result, &status); err != nil {
		return RequestStatus{}, err
	}
	tuple := status.MustObj().MustVec()
	return RequestStatus{
		Fulfilled:      tuple[0].MustIc() == xdr.ScStaticScsTrue,
		LedgerSequence: uint32(tuple[1].MustU32()),
		CloseTime:      time.Unix(int64(tuple[2].MustObj().MustU64()), 0),
	}, nil
}
