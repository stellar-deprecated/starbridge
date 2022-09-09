// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RegisterStellarAssetRequest is an auto generated low-level Go binding around an user-defined struct.
type RegisterStellarAssetRequest struct {
	Decimals uint8
	Name     string
	Symbol   string
}

// SetPausedRequest is an auto generated low-level Go binding around an user-defined struct.
type SetPausedRequest struct {
	Value      uint8
	Nonce      *big.Int
	Expiration *big.Int
}

// WithdrawERC20Request is an auto generated low-level Go binding around an user-defined struct.
type WithdrawERC20Request struct {
	Id         [32]byte
	Expiration *big.Int
	Recipient  common.Address
	Token      common.Address
	Amount     *big.Int
}

// WithdrawETHRequest is an auto generated low-level Go binding around an user-defined struct.
type WithdrawETHRequest struct {
	Id         [32]byte
	Expiration *big.Int
	Recipient  common.Address
	Amount     *big.Int
}

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_minThreshold\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"minThreshold\",\"type\":\"uint8\"}],\"name\":\"RegisterSigners\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"RegisterStellarAsset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"}],\"name\":\"SetPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"}],\"name\":\"depositETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"domainSeparator\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isStellarAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minThreshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"internalType\":\"structRegisterStellarAssetRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8[]\",\"name\":\"indexes\",\"type\":\"uint8[]\"}],\"name\":\"registerStellarAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestID\",\"type\":\"bytes32\"}],\"name\":\"requestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structSetPausedRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8[]\",\"name\":\"indexes\",\"type\":\"uint8[]\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_minThreshold\",\"type\":\"uint8\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8[]\",\"name\":\"indexes\",\"type\":\"uint8[]\"}],\"name\":\"updateSigners\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawERC20Request\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8[]\",\"name\":\"indexes\",\"type\":\"uint8[]\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawETHRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8[]\",\"name\":\"indexes\",\"type\":\"uint8[]\"}],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Bridge *BridgeCaller) DomainSeparator(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "domainSeparator")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Bridge *BridgeSession) DomainSeparator() ([32]byte, error) {
	return _Bridge.Contract.DomainSeparator(&_Bridge.CallOpts)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Bridge *BridgeCallerSession) DomainSeparator() ([32]byte, error) {
	return _Bridge.Contract.DomainSeparator(&_Bridge.CallOpts)
}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_Bridge *BridgeCaller) IsStellarAsset(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "isStellarAsset", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_Bridge *BridgeSession) IsStellarAsset(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.IsStellarAsset(&_Bridge.CallOpts, arg0)
}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_Bridge *BridgeCallerSession) IsStellarAsset(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.IsStellarAsset(&_Bridge.CallOpts, arg0)
}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint8)
func (_Bridge *BridgeCaller) MinThreshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "minThreshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint8)
func (_Bridge *BridgeSession) MinThreshold() (uint8, error) {
	return _Bridge.Contract.MinThreshold(&_Bridge.CallOpts)
}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint8)
func (_Bridge *BridgeCallerSession) MinThreshold() (uint8, error) {
	return _Bridge.Contract.MinThreshold(&_Bridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint8)
func (_Bridge *BridgeCaller) Paused(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint8)
func (_Bridge *BridgeSession) Paused() (uint8, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint8)
func (_Bridge *BridgeCallerSession) Paused() (uint8, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// RequestStatus is a free data retrieval call binding the contract method 0x9b902497.
//
// Solidity: function requestStatus(bytes32 requestID) view returns(bool, uint256)
func (_Bridge *BridgeCaller) RequestStatus(opts *bind.CallOpts, requestID [32]byte) (bool, *big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "requestStatus", requestID)

	if err != nil {
		return *new(bool), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// RequestStatus is a free data retrieval call binding the contract method 0x9b902497.
//
// Solidity: function requestStatus(bytes32 requestID) view returns(bool, uint256)
func (_Bridge *BridgeSession) RequestStatus(requestID [32]byte) (bool, *big.Int, error) {
	return _Bridge.Contract.RequestStatus(&_Bridge.CallOpts, requestID)
}

// RequestStatus is a free data retrieval call binding the contract method 0x9b902497.
//
// Solidity: function requestStatus(bytes32 requestID) view returns(bool, uint256)
func (_Bridge *BridgeCallerSession) RequestStatus(requestID [32]byte) (bool, *big.Int, error) {
	return _Bridge.Contract.RequestStatus(&_Bridge.CallOpts, requestID)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Bridge *BridgeCaller) Signers(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Bridge *BridgeSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _Bridge.Contract.Signers(&_Bridge.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Bridge *BridgeCallerSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _Bridge.Contract.Signers(&_Bridge.CallOpts, arg0)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Bridge *BridgeCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Bridge *BridgeSession) Version() (*big.Int, error) {
	return _Bridge.Contract.Version(&_Bridge.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Bridge *BridgeCallerSession) Version() (*big.Int, error) {
	return _Bridge.Contract.Version(&_Bridge.CallOpts)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x21425ee0.
//
// Solidity: function depositERC20(address token, uint256 destination, uint256 amount) returns()
func (_Bridge *BridgeTransactor) DepositERC20(opts *bind.TransactOpts, token common.Address, destination *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "depositERC20", token, destination, amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x21425ee0.
//
// Solidity: function depositERC20(address token, uint256 destination, uint256 amount) returns()
func (_Bridge *BridgeSession) DepositERC20(token common.Address, destination *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.DepositERC20(&_Bridge.TransactOpts, token, destination, amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x21425ee0.
//
// Solidity: function depositERC20(address token, uint256 destination, uint256 amount) returns()
func (_Bridge *BridgeTransactorSession) DepositERC20(token common.Address, destination *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.DepositERC20(&_Bridge.TransactOpts, token, destination, amount)
}

// DepositETH is a paid mutator transaction binding the contract method 0x5358fbda.
//
// Solidity: function depositETH(uint256 destination) payable returns()
func (_Bridge *BridgeTransactor) DepositETH(opts *bind.TransactOpts, destination *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "depositETH", destination)
}

// DepositETH is a paid mutator transaction binding the contract method 0x5358fbda.
//
// Solidity: function depositETH(uint256 destination) payable returns()
func (_Bridge *BridgeSession) DepositETH(destination *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.DepositETH(&_Bridge.TransactOpts, destination)
}

// DepositETH is a paid mutator transaction binding the contract method 0x5358fbda.
//
// Solidity: function depositETH(uint256 destination) payable returns()
func (_Bridge *BridgeTransactorSession) DepositETH(destination *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.DepositETH(&_Bridge.TransactOpts, destination)
}

// RegisterStellarAsset is a paid mutator transaction binding the contract method 0x410d8d61.
//
// Solidity: function registerStellarAsset((uint8,string,string) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactor) RegisterStellarAsset(opts *bind.TransactOpts, request RegisterStellarAssetRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerStellarAsset", request, signatures, indexes)
}

// RegisterStellarAsset is a paid mutator transaction binding the contract method 0x410d8d61.
//
// Solidity: function registerStellarAsset((uint8,string,string) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeSession) RegisterStellarAsset(request RegisterStellarAssetRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterStellarAsset(&_Bridge.TransactOpts, request, signatures, indexes)
}

// RegisterStellarAsset is a paid mutator transaction binding the contract method 0x410d8d61.
//
// Solidity: function registerStellarAsset((uint8,string,string) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactorSession) RegisterStellarAsset(request RegisterStellarAssetRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterStellarAsset(&_Bridge.TransactOpts, request, signatures, indexes)
}

// SetPaused is a paid mutator transaction binding the contract method 0xfac7d40f.
//
// Solidity: function setPaused((uint8,uint256,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactor) SetPaused(opts *bind.TransactOpts, request SetPausedRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setPaused", request, signatures, indexes)
}

// SetPaused is a paid mutator transaction binding the contract method 0xfac7d40f.
//
// Solidity: function setPaused((uint8,uint256,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeSession) SetPaused(request SetPausedRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.SetPaused(&_Bridge.TransactOpts, request, signatures, indexes)
}

// SetPaused is a paid mutator transaction binding the contract method 0xfac7d40f.
//
// Solidity: function setPaused((uint8,uint256,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactorSession) SetPaused(request SetPausedRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.SetPaused(&_Bridge.TransactOpts, request, signatures, indexes)
}

// UpdateSigners is a paid mutator transaction binding the contract method 0xbd7c5733.
//
// Solidity: function updateSigners(address[] _signers, uint8 _minThreshold, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactor) UpdateSigners(opts *bind.TransactOpts, _signers []common.Address, _minThreshold uint8, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "updateSigners", _signers, _minThreshold, signatures, indexes)
}

// UpdateSigners is a paid mutator transaction binding the contract method 0xbd7c5733.
//
// Solidity: function updateSigners(address[] _signers, uint8 _minThreshold, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeSession) UpdateSigners(_signers []common.Address, _minThreshold uint8, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.UpdateSigners(&_Bridge.TransactOpts, _signers, _minThreshold, signatures, indexes)
}

// UpdateSigners is a paid mutator transaction binding the contract method 0xbd7c5733.
//
// Solidity: function updateSigners(address[] _signers, uint8 _minThreshold, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactorSession) UpdateSigners(_signers []common.Address, _minThreshold uint8, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.UpdateSigners(&_Bridge.TransactOpts, _signers, _minThreshold, signatures, indexes)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x23d64d1c.
//
// Solidity: function withdrawERC20((bytes32,uint256,address,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactor) WithdrawERC20(opts *bind.TransactOpts, request WithdrawERC20Request, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "withdrawERC20", request, signatures, indexes)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x23d64d1c.
//
// Solidity: function withdrawERC20((bytes32,uint256,address,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeSession) WithdrawERC20(request WithdrawERC20Request, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawERC20(&_Bridge.TransactOpts, request, signatures, indexes)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x23d64d1c.
//
// Solidity: function withdrawERC20((bytes32,uint256,address,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactorSession) WithdrawERC20(request WithdrawERC20Request, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawERC20(&_Bridge.TransactOpts, request, signatures, indexes)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x15c0f43d.
//
// Solidity: function withdrawETH((bytes32,uint256,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactor) WithdrawETH(opts *bind.TransactOpts, request WithdrawETHRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "withdrawETH", request, signatures, indexes)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x15c0f43d.
//
// Solidity: function withdrawETH((bytes32,uint256,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeSession) WithdrawETH(request WithdrawETHRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawETH(&_Bridge.TransactOpts, request, signatures, indexes)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x15c0f43d.
//
// Solidity: function withdrawETH((bytes32,uint256,address,uint256) request, bytes[] signatures, uint8[] indexes) returns()
func (_Bridge *BridgeTransactorSession) WithdrawETH(request WithdrawETHRequest, signatures [][]byte, indexes []uint8) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawETH(&_Bridge.TransactOpts, request, signatures, indexes)
}

// BridgeDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Bridge contract.
type BridgeDepositIterator struct {
	Event *BridgeDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeDeposit represents a Deposit event raised by the Bridge contract.
type BridgeDeposit struct {
	Token       common.Address
	Sender      common.Address
	Destination *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) FilterDeposit(opts *bind.FilterOpts) (*BridgeDepositIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &BridgeDepositIterator{contract: _Bridge.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *BridgeDeposit) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeDeposit)
				if err := _Bridge.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) ParseDeposit(log types.Log) (*BridgeDeposit, error) {
	event := new(BridgeDeposit)
	if err := _Bridge.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeRegisterSignersIterator is returned from FilterRegisterSigners and is used to iterate over the raw logs and unpacked data for RegisterSigners events raised by the Bridge contract.
type BridgeRegisterSignersIterator struct {
	Event *BridgeRegisterSigners // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeRegisterSignersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRegisterSigners)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeRegisterSigners)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeRegisterSignersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRegisterSignersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRegisterSigners represents a RegisterSigners event raised by the Bridge contract.
type BridgeRegisterSigners struct {
	Version         *big.Int
	DomainSeparator [32]byte
	Signers         []common.Address
	MinThreshold    uint8
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRegisterSigners is a free log retrieval operation binding the contract event 0x8efb61e94bf9ccfbdcc92531bc3c176ce376859434de722bec5b6e0813df9f2c.
//
// Solidity: event RegisterSigners(uint256 version, bytes32 domainSeparator, address[] signers, uint8 minThreshold)
func (_Bridge *BridgeFilterer) FilterRegisterSigners(opts *bind.FilterOpts) (*BridgeRegisterSignersIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RegisterSigners")
	if err != nil {
		return nil, err
	}
	return &BridgeRegisterSignersIterator{contract: _Bridge.contract, event: "RegisterSigners", logs: logs, sub: sub}, nil
}

// WatchRegisterSigners is a free log subscription operation binding the contract event 0x8efb61e94bf9ccfbdcc92531bc3c176ce376859434de722bec5b6e0813df9f2c.
//
// Solidity: event RegisterSigners(uint256 version, bytes32 domainSeparator, address[] signers, uint8 minThreshold)
func (_Bridge *BridgeFilterer) WatchRegisterSigners(opts *bind.WatchOpts, sink chan<- *BridgeRegisterSigners) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RegisterSigners")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRegisterSigners)
				if err := _Bridge.contract.UnpackLog(event, "RegisterSigners", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisterSigners is a log parse operation binding the contract event 0x8efb61e94bf9ccfbdcc92531bc3c176ce376859434de722bec5b6e0813df9f2c.
//
// Solidity: event RegisterSigners(uint256 version, bytes32 domainSeparator, address[] signers, uint8 minThreshold)
func (_Bridge *BridgeFilterer) ParseRegisterSigners(log types.Log) (*BridgeRegisterSigners, error) {
	event := new(BridgeRegisterSigners)
	if err := _Bridge.contract.UnpackLog(event, "RegisterSigners", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeRegisterStellarAssetIterator is returned from FilterRegisterStellarAsset and is used to iterate over the raw logs and unpacked data for RegisterStellarAsset events raised by the Bridge contract.
type BridgeRegisterStellarAssetIterator struct {
	Event *BridgeRegisterStellarAsset // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeRegisterStellarAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRegisterStellarAsset)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeRegisterStellarAsset)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeRegisterStellarAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRegisterStellarAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRegisterStellarAsset represents a RegisterStellarAsset event raised by the Bridge contract.
type BridgeRegisterStellarAsset struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRegisterStellarAsset is a free log retrieval operation binding the contract event 0xa4dfd59983ab499bfeddc2f42db87b2c7e274c8c227df1b705c0139d9f7ca3f3.
//
// Solidity: event RegisterStellarAsset(address asset)
func (_Bridge *BridgeFilterer) FilterRegisterStellarAsset(opts *bind.FilterOpts) (*BridgeRegisterStellarAssetIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RegisterStellarAsset")
	if err != nil {
		return nil, err
	}
	return &BridgeRegisterStellarAssetIterator{contract: _Bridge.contract, event: "RegisterStellarAsset", logs: logs, sub: sub}, nil
}

// WatchRegisterStellarAsset is a free log subscription operation binding the contract event 0xa4dfd59983ab499bfeddc2f42db87b2c7e274c8c227df1b705c0139d9f7ca3f3.
//
// Solidity: event RegisterStellarAsset(address asset)
func (_Bridge *BridgeFilterer) WatchRegisterStellarAsset(opts *bind.WatchOpts, sink chan<- *BridgeRegisterStellarAsset) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RegisterStellarAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRegisterStellarAsset)
				if err := _Bridge.contract.UnpackLog(event, "RegisterStellarAsset", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisterStellarAsset is a log parse operation binding the contract event 0xa4dfd59983ab499bfeddc2f42db87b2c7e274c8c227df1b705c0139d9f7ca3f3.
//
// Solidity: event RegisterStellarAsset(address asset)
func (_Bridge *BridgeFilterer) ParseRegisterStellarAsset(log types.Log) (*BridgeRegisterStellarAsset, error) {
	event := new(BridgeRegisterStellarAsset)
	if err := _Bridge.contract.UnpackLog(event, "RegisterStellarAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeSetPausedIterator is returned from FilterSetPaused and is used to iterate over the raw logs and unpacked data for SetPaused events raised by the Bridge contract.
type BridgeSetPausedIterator struct {
	Event *BridgeSetPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeSetPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeSetPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeSetPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeSetPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeSetPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeSetPaused represents a SetPaused event raised by the Bridge contract.
type BridgeSetPaused struct {
	Value uint8
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSetPaused is a free log retrieval operation binding the contract event 0x9a0d3d5e378f6c06261926f8d2c566764cc506766d4cf6535cb64d9f831fa310.
//
// Solidity: event SetPaused(uint8 value)
func (_Bridge *BridgeFilterer) FilterSetPaused(opts *bind.FilterOpts) (*BridgeSetPausedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "SetPaused")
	if err != nil {
		return nil, err
	}
	return &BridgeSetPausedIterator{contract: _Bridge.contract, event: "SetPaused", logs: logs, sub: sub}, nil
}

// WatchSetPaused is a free log subscription operation binding the contract event 0x9a0d3d5e378f6c06261926f8d2c566764cc506766d4cf6535cb64d9f831fa310.
//
// Solidity: event SetPaused(uint8 value)
func (_Bridge *BridgeFilterer) WatchSetPaused(opts *bind.WatchOpts, sink chan<- *BridgeSetPaused) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "SetPaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeSetPaused)
				if err := _Bridge.contract.UnpackLog(event, "SetPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetPaused is a log parse operation binding the contract event 0x9a0d3d5e378f6c06261926f8d2c566764cc506766d4cf6535cb64d9f831fa310.
//
// Solidity: event SetPaused(uint8 value)
func (_Bridge *BridgeFilterer) ParseSetPaused(log types.Log) (*BridgeSetPaused, error) {
	event := new(BridgeSetPaused)
	if err := _Bridge.contract.UnpackLog(event, "SetPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Bridge contract.
type BridgeWithdrawIterator struct {
	Event *BridgeWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeWithdraw represents a Withdraw event raised by the Bridge contract.
type BridgeWithdraw struct {
	Id        [32]byte
	Token     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xcd6ac346191b4b7531743e58f243dd4d350a52a9186641c1e5eac22b95aaedbe.
//
// Solidity: event Withdraw(bytes32 id, address token, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) FilterWithdraw(opts *bind.FilterOpts) (*BridgeWithdrawIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &BridgeWithdrawIterator{contract: _Bridge.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xcd6ac346191b4b7531743e58f243dd4d350a52a9186641c1e5eac22b95aaedbe.
//
// Solidity: event Withdraw(bytes32 id, address token, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *BridgeWithdraw) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeWithdraw)
				if err := _Bridge.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0xcd6ac346191b4b7531743e58f243dd4d350a52a9186641c1e5eac22b95aaedbe.
//
// Solidity: event Withdraw(bytes32 id, address token, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) ParseWithdraw(log types.Log) (*BridgeWithdraw, error) {
	event := new(BridgeWithdraw)
	if err := _Bridge.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
