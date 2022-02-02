// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package soliditybridge

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

// AuthInternalMintStellarAssetRequest is an auto generated low-level Go binding around an user-defined struct.
type AuthInternalMintStellarAssetRequest struct {
	Nonce     *big.Int
	Recipient common.Address
	Amount    *big.Int
	Decimals  uint8
}

// MintStellarAssetRequest is an auto generated low-level Go binding around an user-defined struct.
type MintStellarAssetRequest struct {
	Nonce     *big.Int
	Recipient common.Address
	Amount    *big.Int
	Decimals  uint8
	Name      string
	Symbol    string
}

// WithdrawERC20Request is an auto generated low-level Go binding around an user-defined struct.
type WithdrawERC20Request struct {
	Nonce     *big.Int
	Recipient common.Address
	Token     common.Address
	Amount    *big.Int
}

// WithdrawETHRequest is an auto generated low-level Go binding around an user-defined struct.
type WithdrawETHRequest struct {
	Nonce     *big.Int
	Recipient common.Address
	Amount    *big.Int
}

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnStellarAsset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"CreateStellarAsset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositERC20\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositETH\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintStellarAsset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"RegisterSigners\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawERC20\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawETH\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINT_STELLAR_ASSET_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW_ERC20_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW_ETH_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"destination\",\"type\":\"uint256\"}],\"name\":\"depositETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"internalType\":\"structAuth.InternalMintStellarAssetRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"nameHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"symbolHash\",\"type\":\"bytes32\"}],\"name\":\"hashMintStellarAssetRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawERC20Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"hashWithdrawERC20Request\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawETHRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"hashWithdrawETHRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isStellarAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"internalType\":\"structMintStellarAssetRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"mintStellarAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nameHashToAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawERC20Request\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawETHRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Bridge *BridgeCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Bridge *BridgeSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Bridge.Contract.DOMAINSEPARATOR(&_Bridge.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Bridge *BridgeCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Bridge.Contract.DOMAINSEPARATOR(&_Bridge.CallOpts)
}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCaller) MINTSTELLARASSETTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "MINT_STELLAR_ASSET_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeSession) MINTSTELLARASSETTYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.MINTSTELLARASSETTYPEHASH(&_Bridge.CallOpts)
}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCallerSession) MINTSTELLARASSETTYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.MINTSTELLARASSETTYPEHASH(&_Bridge.CallOpts)
}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCaller) WITHDRAWERC20TYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "WITHDRAW_ERC20_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeSession) WITHDRAWERC20TYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.WITHDRAWERC20TYPEHASH(&_Bridge.CallOpts)
}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCallerSession) WITHDRAWERC20TYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.WITHDRAWERC20TYPEHASH(&_Bridge.CallOpts)
}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCaller) WITHDRAWETHTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "WITHDRAW_ETH_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeSession) WITHDRAWETHTYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.WITHDRAWETHTYPEHASH(&_Bridge.CallOpts)
}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Bridge *BridgeCallerSession) WITHDRAWETHTYPEHASH() ([32]byte, error) {
	return _Bridge.Contract.WITHDRAWETHTYPEHASH(&_Bridge.CallOpts)
}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Bridge *BridgeCaller) HashMintStellarAssetRequest(opts *bind.CallOpts, request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "hashMintStellarAssetRequest", request, nameHash, symbolHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Bridge *BridgeSession) HashMintStellarAssetRequest(request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	return _Bridge.Contract.HashMintStellarAssetRequest(&_Bridge.CallOpts, request, nameHash, symbolHash)
}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Bridge *BridgeCallerSession) HashMintStellarAssetRequest(request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	return _Bridge.Contract.HashMintStellarAssetRequest(&_Bridge.CallOpts, request, nameHash, symbolHash)
}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeCaller) HashWithdrawERC20Request(opts *bind.CallOpts, request WithdrawERC20Request) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "hashWithdrawERC20Request", request)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeSession) HashWithdrawERC20Request(request WithdrawERC20Request) ([32]byte, error) {
	return _Bridge.Contract.HashWithdrawERC20Request(&_Bridge.CallOpts, request)
}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeCallerSession) HashWithdrawERC20Request(request WithdrawERC20Request) ([32]byte, error) {
	return _Bridge.Contract.HashWithdrawERC20Request(&_Bridge.CallOpts, request)
}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeCaller) HashWithdrawETHRequest(opts *bind.CallOpts, request WithdrawETHRequest) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "hashWithdrawETHRequest", request)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeSession) HashWithdrawETHRequest(request WithdrawETHRequest) ([32]byte, error) {
	return _Bridge.Contract.HashWithdrawETHRequest(&_Bridge.CallOpts, request)
}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Bridge *BridgeCallerSession) HashWithdrawETHRequest(request WithdrawETHRequest) ([32]byte, error) {
	return _Bridge.Contract.HashWithdrawETHRequest(&_Bridge.CallOpts, request)
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

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_Bridge *BridgeCaller) NameHashToAsset(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "nameHashToAsset", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_Bridge *BridgeSession) NameHashToAsset(arg0 [32]byte) (common.Address, error) {
	return _Bridge.Contract.NameHashToAsset(&_Bridge.CallOpts, arg0)
}

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_Bridge *BridgeCallerSession) NameHashToAsset(arg0 [32]byte) (common.Address, error) {
	return _Bridge.Contract.NameHashToAsset(&_Bridge.CallOpts, arg0)
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

// MintStellarAsset is a paid mutator transaction binding the contract method 0x1738c0dc.
//
// Solidity: function mintStellarAsset((uint256,address,uint256,uint8,string,string) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactor) MintStellarAsset(opts *bind.TransactOpts, request MintStellarAssetRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "mintStellarAsset", request, signatures)
}

// MintStellarAsset is a paid mutator transaction binding the contract method 0x1738c0dc.
//
// Solidity: function mintStellarAsset((uint256,address,uint256,uint8,string,string) request, bytes[] signatures) returns()
func (_Bridge *BridgeSession) MintStellarAsset(request MintStellarAssetRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.MintStellarAsset(&_Bridge.TransactOpts, request, signatures)
}

// MintStellarAsset is a paid mutator transaction binding the contract method 0x1738c0dc.
//
// Solidity: function mintStellarAsset((uint256,address,uint256,uint8,string,string) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactorSession) MintStellarAsset(request MintStellarAssetRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.MintStellarAsset(&_Bridge.TransactOpts, request, signatures)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa29f7f08.
//
// Solidity: function withdrawERC20((uint256,address,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactor) WithdrawERC20(opts *bind.TransactOpts, request WithdrawERC20Request, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "withdrawERC20", request, signatures)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa29f7f08.
//
// Solidity: function withdrawERC20((uint256,address,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeSession) WithdrawERC20(request WithdrawERC20Request, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawERC20(&_Bridge.TransactOpts, request, signatures)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa29f7f08.
//
// Solidity: function withdrawERC20((uint256,address,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactorSession) WithdrawERC20(request WithdrawERC20Request, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawERC20(&_Bridge.TransactOpts, request, signatures)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x0fe1fcbb.
//
// Solidity: function withdrawETH((uint256,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactor) WithdrawETH(opts *bind.TransactOpts, request WithdrawETHRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "withdrawETH", request, signatures)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x0fe1fcbb.
//
// Solidity: function withdrawETH((uint256,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeSession) WithdrawETH(request WithdrawETHRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawETH(&_Bridge.TransactOpts, request, signatures)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x0fe1fcbb.
//
// Solidity: function withdrawETH((uint256,address,uint256) request, bytes[] signatures) returns()
func (_Bridge *BridgeTransactorSession) WithdrawETH(request WithdrawETHRequest, signatures [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.WithdrawETH(&_Bridge.TransactOpts, request, signatures)
}

// BridgeBurnStellarAssetIterator is returned from FilterBurnStellarAsset and is used to iterate over the raw logs and unpacked data for BurnStellarAsset events raised by the Bridge contract.
type BridgeBurnStellarAssetIterator struct {
	Event *BridgeBurnStellarAsset // Event containing the contract specifics and raw log

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
func (it *BridgeBurnStellarAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeBurnStellarAsset)
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
		it.Event = new(BridgeBurnStellarAsset)
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
func (it *BridgeBurnStellarAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeBurnStellarAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeBurnStellarAsset represents a BurnStellarAsset event raised by the Bridge contract.
type BridgeBurnStellarAsset struct {
	Asset  common.Address
	Holder common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurnStellarAsset is a free log retrieval operation binding the contract event 0xb832502e867c209b6e037c31c0a7d3be7b47fb3c54e49fc926d3b691f4a32fc6.
//
// Solidity: event BurnStellarAsset(address asset, address holder, uint256 amount)
func (_Bridge *BridgeFilterer) FilterBurnStellarAsset(opts *bind.FilterOpts) (*BridgeBurnStellarAssetIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "BurnStellarAsset")
	if err != nil {
		return nil, err
	}
	return &BridgeBurnStellarAssetIterator{contract: _Bridge.contract, event: "BurnStellarAsset", logs: logs, sub: sub}, nil
}

// WatchBurnStellarAsset is a free log subscription operation binding the contract event 0xb832502e867c209b6e037c31c0a7d3be7b47fb3c54e49fc926d3b691f4a32fc6.
//
// Solidity: event BurnStellarAsset(address asset, address holder, uint256 amount)
func (_Bridge *BridgeFilterer) WatchBurnStellarAsset(opts *bind.WatchOpts, sink chan<- *BridgeBurnStellarAsset) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "BurnStellarAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeBurnStellarAsset)
				if err := _Bridge.contract.UnpackLog(event, "BurnStellarAsset", log); err != nil {
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

// ParseBurnStellarAsset is a log parse operation binding the contract event 0xb832502e867c209b6e037c31c0a7d3be7b47fb3c54e49fc926d3b691f4a32fc6.
//
// Solidity: event BurnStellarAsset(address asset, address holder, uint256 amount)
func (_Bridge *BridgeFilterer) ParseBurnStellarAsset(log types.Log) (*BridgeBurnStellarAsset, error) {
	event := new(BridgeBurnStellarAsset)
	if err := _Bridge.contract.UnpackLog(event, "BurnStellarAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeCreateStellarAssetIterator is returned from FilterCreateStellarAsset and is used to iterate over the raw logs and unpacked data for CreateStellarAsset events raised by the Bridge contract.
type BridgeCreateStellarAssetIterator struct {
	Event *BridgeCreateStellarAsset // Event containing the contract specifics and raw log

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
func (it *BridgeCreateStellarAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeCreateStellarAsset)
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
		it.Event = new(BridgeCreateStellarAsset)
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
func (it *BridgeCreateStellarAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeCreateStellarAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeCreateStellarAsset represents a CreateStellarAsset event raised by the Bridge contract.
type BridgeCreateStellarAsset struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCreateStellarAsset is a free log retrieval operation binding the contract event 0x103935973961d03366259eb4fc092bbc84b0812cba6d190f70d03f5f7bbf4cc2.
//
// Solidity: event CreateStellarAsset(address asset)
func (_Bridge *BridgeFilterer) FilterCreateStellarAsset(opts *bind.FilterOpts) (*BridgeCreateStellarAssetIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "CreateStellarAsset")
	if err != nil {
		return nil, err
	}
	return &BridgeCreateStellarAssetIterator{contract: _Bridge.contract, event: "CreateStellarAsset", logs: logs, sub: sub}, nil
}

// WatchCreateStellarAsset is a free log subscription operation binding the contract event 0x103935973961d03366259eb4fc092bbc84b0812cba6d190f70d03f5f7bbf4cc2.
//
// Solidity: event CreateStellarAsset(address asset)
func (_Bridge *BridgeFilterer) WatchCreateStellarAsset(opts *bind.WatchOpts, sink chan<- *BridgeCreateStellarAsset) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "CreateStellarAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeCreateStellarAsset)
				if err := _Bridge.contract.UnpackLog(event, "CreateStellarAsset", log); err != nil {
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

// ParseCreateStellarAsset is a log parse operation binding the contract event 0x103935973961d03366259eb4fc092bbc84b0812cba6d190f70d03f5f7bbf4cc2.
//
// Solidity: event CreateStellarAsset(address asset)
func (_Bridge *BridgeFilterer) ParseCreateStellarAsset(log types.Log) (*BridgeCreateStellarAsset, error) {
	event := new(BridgeCreateStellarAsset)
	if err := _Bridge.contract.UnpackLog(event, "CreateStellarAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeDepositERC20Iterator is returned from FilterDepositERC20 and is used to iterate over the raw logs and unpacked data for DepositERC20 events raised by the Bridge contract.
type BridgeDepositERC20Iterator struct {
	Event *BridgeDepositERC20 // Event containing the contract specifics and raw log

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
func (it *BridgeDepositERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeDepositERC20)
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
		it.Event = new(BridgeDepositERC20)
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
func (it *BridgeDepositERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeDepositERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeDepositERC20 represents a DepositERC20 event raised by the Bridge contract.
type BridgeDepositERC20 struct {
	Token       common.Address
	Sender      common.Address
	Destination *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDepositERC20 is a free log retrieval operation binding the contract event 0x228adef123a7fe6b727889b4d24b39b39de663428214fd757c6c501fb98d8494.
//
// Solidity: event DepositERC20(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) FilterDepositERC20(opts *bind.FilterOpts) (*BridgeDepositERC20Iterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "DepositERC20")
	if err != nil {
		return nil, err
	}
	return &BridgeDepositERC20Iterator{contract: _Bridge.contract, event: "DepositERC20", logs: logs, sub: sub}, nil
}

// WatchDepositERC20 is a free log subscription operation binding the contract event 0x228adef123a7fe6b727889b4d24b39b39de663428214fd757c6c501fb98d8494.
//
// Solidity: event DepositERC20(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) WatchDepositERC20(opts *bind.WatchOpts, sink chan<- *BridgeDepositERC20) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "DepositERC20")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeDepositERC20)
				if err := _Bridge.contract.UnpackLog(event, "DepositERC20", log); err != nil {
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

// ParseDepositERC20 is a log parse operation binding the contract event 0x228adef123a7fe6b727889b4d24b39b39de663428214fd757c6c501fb98d8494.
//
// Solidity: event DepositERC20(address token, address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) ParseDepositERC20(log types.Log) (*BridgeDepositERC20, error) {
	event := new(BridgeDepositERC20)
	if err := _Bridge.contract.UnpackLog(event, "DepositERC20", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeDepositETHIterator is returned from FilterDepositETH and is used to iterate over the raw logs and unpacked data for DepositETH events raised by the Bridge contract.
type BridgeDepositETHIterator struct {
	Event *BridgeDepositETH // Event containing the contract specifics and raw log

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
func (it *BridgeDepositETHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeDepositETH)
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
		it.Event = new(BridgeDepositETH)
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
func (it *BridgeDepositETHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeDepositETHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeDepositETH represents a DepositETH event raised by the Bridge contract.
type BridgeDepositETH struct {
	Sender      common.Address
	Destination *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDepositETH is a free log retrieval operation binding the contract event 0x57e8e547a3ef8d890c570ca885b0a8c441be3070e36b7ad4c7d6b9d9316ff2ce.
//
// Solidity: event DepositETH(address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) FilterDepositETH(opts *bind.FilterOpts) (*BridgeDepositETHIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "DepositETH")
	if err != nil {
		return nil, err
	}
	return &BridgeDepositETHIterator{contract: _Bridge.contract, event: "DepositETH", logs: logs, sub: sub}, nil
}

// WatchDepositETH is a free log subscription operation binding the contract event 0x57e8e547a3ef8d890c570ca885b0a8c441be3070e36b7ad4c7d6b9d9316ff2ce.
//
// Solidity: event DepositETH(address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) WatchDepositETH(opts *bind.WatchOpts, sink chan<- *BridgeDepositETH) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "DepositETH")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeDepositETH)
				if err := _Bridge.contract.UnpackLog(event, "DepositETH", log); err != nil {
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

// ParseDepositETH is a log parse operation binding the contract event 0x57e8e547a3ef8d890c570ca885b0a8c441be3070e36b7ad4c7d6b9d9316ff2ce.
//
// Solidity: event DepositETH(address sender, uint256 destination, uint256 amount)
func (_Bridge *BridgeFilterer) ParseDepositETH(log types.Log) (*BridgeDepositETH, error) {
	event := new(BridgeDepositETH)
	if err := _Bridge.contract.UnpackLog(event, "DepositETH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeMintStellarAssetIterator is returned from FilterMintStellarAsset and is used to iterate over the raw logs and unpacked data for MintStellarAsset events raised by the Bridge contract.
type BridgeMintStellarAssetIterator struct {
	Event *BridgeMintStellarAsset // Event containing the contract specifics and raw log

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
func (it *BridgeMintStellarAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeMintStellarAsset)
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
		it.Event = new(BridgeMintStellarAsset)
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
func (it *BridgeMintStellarAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeMintStellarAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeMintStellarAsset represents a MintStellarAsset event raised by the Bridge contract.
type BridgeMintStellarAsset struct {
	Nonce     *big.Int
	Asset     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMintStellarAsset is a free log retrieval operation binding the contract event 0x48728d6a6a9f2e1ca650ecc330e2ba43cbb36520e76c9be9b385c01605cc32d5.
//
// Solidity: event MintStellarAsset(uint256 nonce, address asset, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) FilterMintStellarAsset(opts *bind.FilterOpts) (*BridgeMintStellarAssetIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "MintStellarAsset")
	if err != nil {
		return nil, err
	}
	return &BridgeMintStellarAssetIterator{contract: _Bridge.contract, event: "MintStellarAsset", logs: logs, sub: sub}, nil
}

// WatchMintStellarAsset is a free log subscription operation binding the contract event 0x48728d6a6a9f2e1ca650ecc330e2ba43cbb36520e76c9be9b385c01605cc32d5.
//
// Solidity: event MintStellarAsset(uint256 nonce, address asset, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) WatchMintStellarAsset(opts *bind.WatchOpts, sink chan<- *BridgeMintStellarAsset) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "MintStellarAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeMintStellarAsset)
				if err := _Bridge.contract.UnpackLog(event, "MintStellarAsset", log); err != nil {
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

// ParseMintStellarAsset is a log parse operation binding the contract event 0x48728d6a6a9f2e1ca650ecc330e2ba43cbb36520e76c9be9b385c01605cc32d5.
//
// Solidity: event MintStellarAsset(uint256 nonce, address asset, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) ParseMintStellarAsset(log types.Log) (*BridgeMintStellarAsset, error) {
	event := new(BridgeMintStellarAsset)
	if err := _Bridge.contract.UnpackLog(event, "MintStellarAsset", log); err != nil {
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
	Signers []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRegisterSigners is a free log retrieval operation binding the contract event 0x17ef76377f15d668f61d70a53ea5efbe0a4417c081652ede13ec91f1cee880b0.
//
// Solidity: event RegisterSigners(address[] signers)
func (_Bridge *BridgeFilterer) FilterRegisterSigners(opts *bind.FilterOpts) (*BridgeRegisterSignersIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RegisterSigners")
	if err != nil {
		return nil, err
	}
	return &BridgeRegisterSignersIterator{contract: _Bridge.contract, event: "RegisterSigners", logs: logs, sub: sub}, nil
}

// WatchRegisterSigners is a free log subscription operation binding the contract event 0x17ef76377f15d668f61d70a53ea5efbe0a4417c081652ede13ec91f1cee880b0.
//
// Solidity: event RegisterSigners(address[] signers)
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

// ParseRegisterSigners is a log parse operation binding the contract event 0x17ef76377f15d668f61d70a53ea5efbe0a4417c081652ede13ec91f1cee880b0.
//
// Solidity: event RegisterSigners(address[] signers)
func (_Bridge *BridgeFilterer) ParseRegisterSigners(log types.Log) (*BridgeRegisterSigners, error) {
	event := new(BridgeRegisterSigners)
	if err := _Bridge.contract.UnpackLog(event, "RegisterSigners", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeWithdrawERC20Iterator is returned from FilterWithdrawERC20 and is used to iterate over the raw logs and unpacked data for WithdrawERC20 events raised by the Bridge contract.
type BridgeWithdrawERC20Iterator struct {
	Event *BridgeWithdrawERC20 // Event containing the contract specifics and raw log

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
func (it *BridgeWithdrawERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeWithdrawERC20)
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
		it.Event = new(BridgeWithdrawERC20)
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
func (it *BridgeWithdrawERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeWithdrawERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeWithdrawERC20 represents a WithdrawERC20 event raised by the Bridge contract.
type BridgeWithdrawERC20 struct {
	Nonce     *big.Int
	Recipient common.Address
	Token     common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawERC20 is a free log retrieval operation binding the contract event 0xf8eff9ed148da048feccebca2617170aec98bea5840b0b0be8a289ec03c375fd.
//
// Solidity: event WithdrawERC20(uint256 nonce, address recipient, address token, uint256 amount)
func (_Bridge *BridgeFilterer) FilterWithdrawERC20(opts *bind.FilterOpts) (*BridgeWithdrawERC20Iterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "WithdrawERC20")
	if err != nil {
		return nil, err
	}
	return &BridgeWithdrawERC20Iterator{contract: _Bridge.contract, event: "WithdrawERC20", logs: logs, sub: sub}, nil
}

// WatchWithdrawERC20 is a free log subscription operation binding the contract event 0xf8eff9ed148da048feccebca2617170aec98bea5840b0b0be8a289ec03c375fd.
//
// Solidity: event WithdrawERC20(uint256 nonce, address recipient, address token, uint256 amount)
func (_Bridge *BridgeFilterer) WatchWithdrawERC20(opts *bind.WatchOpts, sink chan<- *BridgeWithdrawERC20) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "WithdrawERC20")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeWithdrawERC20)
				if err := _Bridge.contract.UnpackLog(event, "WithdrawERC20", log); err != nil {
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

// ParseWithdrawERC20 is a log parse operation binding the contract event 0xf8eff9ed148da048feccebca2617170aec98bea5840b0b0be8a289ec03c375fd.
//
// Solidity: event WithdrawERC20(uint256 nonce, address recipient, address token, uint256 amount)
func (_Bridge *BridgeFilterer) ParseWithdrawERC20(log types.Log) (*BridgeWithdrawERC20, error) {
	event := new(BridgeWithdrawERC20)
	if err := _Bridge.contract.UnpackLog(event, "WithdrawERC20", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeWithdrawETHIterator is returned from FilterWithdrawETH and is used to iterate over the raw logs and unpacked data for WithdrawETH events raised by the Bridge contract.
type BridgeWithdrawETHIterator struct {
	Event *BridgeWithdrawETH // Event containing the contract specifics and raw log

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
func (it *BridgeWithdrawETHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeWithdrawETH)
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
		it.Event = new(BridgeWithdrawETH)
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
func (it *BridgeWithdrawETHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeWithdrawETHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeWithdrawETH represents a WithdrawETH event raised by the Bridge contract.
type BridgeWithdrawETH struct {
	Nonce     *big.Int
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawETH is a free log retrieval operation binding the contract event 0x4d7574efc376e1ee5e6eacc9b23ed30ae9f40acfddb028be515fb99a14e2290c.
//
// Solidity: event WithdrawETH(uint256 nonce, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) FilterWithdrawETH(opts *bind.FilterOpts) (*BridgeWithdrawETHIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "WithdrawETH")
	if err != nil {
		return nil, err
	}
	return &BridgeWithdrawETHIterator{contract: _Bridge.contract, event: "WithdrawETH", logs: logs, sub: sub}, nil
}

// WatchWithdrawETH is a free log subscription operation binding the contract event 0x4d7574efc376e1ee5e6eacc9b23ed30ae9f40acfddb028be515fb99a14e2290c.
//
// Solidity: event WithdrawETH(uint256 nonce, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) WatchWithdrawETH(opts *bind.WatchOpts, sink chan<- *BridgeWithdrawETH) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "WithdrawETH")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeWithdrawETH)
				if err := _Bridge.contract.UnpackLog(event, "WithdrawETH", log); err != nil {
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

// ParseWithdrawETH is a log parse operation binding the contract event 0x4d7574efc376e1ee5e6eacc9b23ed30ae9f40acfddb028be515fb99a14e2290c.
//
// Solidity: event WithdrawETH(uint256 nonce, address recipient, uint256 amount)
func (_Bridge *BridgeFilterer) ParseWithdrawETH(log types.Log) (*BridgeWithdrawETH, error) {
	event := new(BridgeWithdrawETH)
	if err := _Bridge.contract.UnpackLog(event, "WithdrawETH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
