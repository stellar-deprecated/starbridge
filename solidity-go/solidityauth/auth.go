// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidityauth

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

// AuthMetaData contains all meta data concerning the Auth contract.
var AuthMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"RegisterSigners\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINT_STELLAR_ASSET_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW_ERC20_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW_ETH_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"internalType\":\"structAuth.InternalMintStellarAssetRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"nameHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"symbolHash\",\"type\":\"bytes32\"}],\"name\":\"hashMintStellarAssetRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawERC20Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"hashWithdrawERC20Request\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structWithdrawETHRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"hashWithdrawETHRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AuthABI is the input ABI used to generate the binding from.
// Deprecated: Use AuthMetaData.ABI instead.
var AuthABI = AuthMetaData.ABI

// Auth is an auto generated Go binding around an Ethereum contract.
type Auth struct {
	AuthCaller     // Read-only binding to the contract
	AuthTransactor // Write-only binding to the contract
	AuthFilterer   // Log filterer for contract events
}

// AuthCaller is an auto generated read-only Go binding around an Ethereum contract.
type AuthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AuthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AuthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AuthSession struct {
	Contract     *Auth             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AuthCallerSession struct {
	Contract *AuthCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AuthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AuthTransactorSession struct {
	Contract     *AuthTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuthRaw is an auto generated low-level Go binding around an Ethereum contract.
type AuthRaw struct {
	Contract *Auth // Generic contract binding to access the raw methods on
}

// AuthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AuthCallerRaw struct {
	Contract *AuthCaller // Generic read-only contract binding to access the raw methods on
}

// AuthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AuthTransactorRaw struct {
	Contract *AuthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAuth creates a new instance of Auth, bound to a specific deployed contract.
func NewAuth(address common.Address, backend bind.ContractBackend) (*Auth, error) {
	contract, err := bindAuth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Auth{AuthCaller: AuthCaller{contract: contract}, AuthTransactor: AuthTransactor{contract: contract}, AuthFilterer: AuthFilterer{contract: contract}}, nil
}

// NewAuthCaller creates a new read-only instance of Auth, bound to a specific deployed contract.
func NewAuthCaller(address common.Address, caller bind.ContractCaller) (*AuthCaller, error) {
	contract, err := bindAuth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuthCaller{contract: contract}, nil
}

// NewAuthTransactor creates a new write-only instance of Auth, bound to a specific deployed contract.
func NewAuthTransactor(address common.Address, transactor bind.ContractTransactor) (*AuthTransactor, error) {
	contract, err := bindAuth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuthTransactor{contract: contract}, nil
}

// NewAuthFilterer creates a new log filterer instance of Auth, bound to a specific deployed contract.
func NewAuthFilterer(address common.Address, filterer bind.ContractFilterer) (*AuthFilterer, error) {
	contract, err := bindAuth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuthFilterer{contract: contract}, nil
}

// bindAuth binds a generic wrapper to an already deployed contract.
func bindAuth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auth *AuthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auth.Contract.AuthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auth *AuthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auth.Contract.AuthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auth *AuthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auth.Contract.AuthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auth *AuthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auth *AuthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auth *AuthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auth.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Auth *AuthCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Auth *AuthSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Auth.Contract.DOMAINSEPARATOR(&_Auth.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Auth *AuthCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Auth.Contract.DOMAINSEPARATOR(&_Auth.CallOpts)
}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCaller) MINTSTELLARASSETTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "MINT_STELLAR_ASSET_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Auth *AuthSession) MINTSTELLARASSETTYPEHASH() ([32]byte, error) {
	return _Auth.Contract.MINTSTELLARASSETTYPEHASH(&_Auth.CallOpts)
}

// MINTSTELLARASSETTYPEHASH is a free data retrieval call binding the contract method 0x853d5211.
//
// Solidity: function MINT_STELLAR_ASSET_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCallerSession) MINTSTELLARASSETTYPEHASH() ([32]byte, error) {
	return _Auth.Contract.MINTSTELLARASSETTYPEHASH(&_Auth.CallOpts)
}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCaller) WITHDRAWERC20TYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "WITHDRAW_ERC20_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Auth *AuthSession) WITHDRAWERC20TYPEHASH() ([32]byte, error) {
	return _Auth.Contract.WITHDRAWERC20TYPEHASH(&_Auth.CallOpts)
}

// WITHDRAWERC20TYPEHASH is a free data retrieval call binding the contract method 0x2e9818a5.
//
// Solidity: function WITHDRAW_ERC20_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCallerSession) WITHDRAWERC20TYPEHASH() ([32]byte, error) {
	return _Auth.Contract.WITHDRAWERC20TYPEHASH(&_Auth.CallOpts)
}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCaller) WITHDRAWETHTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "WITHDRAW_ETH_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Auth *AuthSession) WITHDRAWETHTYPEHASH() ([32]byte, error) {
	return _Auth.Contract.WITHDRAWETHTYPEHASH(&_Auth.CallOpts)
}

// WITHDRAWETHTYPEHASH is a free data retrieval call binding the contract method 0x3a2f0cc9.
//
// Solidity: function WITHDRAW_ETH_TYPEHASH() view returns(bytes32)
func (_Auth *AuthCallerSession) WITHDRAWETHTYPEHASH() ([32]byte, error) {
	return _Auth.Contract.WITHDRAWETHTYPEHASH(&_Auth.CallOpts)
}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Auth *AuthCaller) HashMintStellarAssetRequest(opts *bind.CallOpts, request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "hashMintStellarAssetRequest", request, nameHash, symbolHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Auth *AuthSession) HashMintStellarAssetRequest(request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	return _Auth.Contract.HashMintStellarAssetRequest(&_Auth.CallOpts, request, nameHash, symbolHash)
}

// HashMintStellarAssetRequest is a free data retrieval call binding the contract method 0xfb388fcf.
//
// Solidity: function hashMintStellarAssetRequest((uint256,address,uint256,uint8) request, bytes32 nameHash, bytes32 symbolHash) view returns(bytes32)
func (_Auth *AuthCallerSession) HashMintStellarAssetRequest(request AuthInternalMintStellarAssetRequest, nameHash [32]byte, symbolHash [32]byte) ([32]byte, error) {
	return _Auth.Contract.HashMintStellarAssetRequest(&_Auth.CallOpts, request, nameHash, symbolHash)
}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Auth *AuthCaller) HashWithdrawERC20Request(opts *bind.CallOpts, request WithdrawERC20Request) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "hashWithdrawERC20Request", request)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Auth *AuthSession) HashWithdrawERC20Request(request WithdrawERC20Request) ([32]byte, error) {
	return _Auth.Contract.HashWithdrawERC20Request(&_Auth.CallOpts, request)
}

// HashWithdrawERC20Request is a free data retrieval call binding the contract method 0x2d7746c8.
//
// Solidity: function hashWithdrawERC20Request((uint256,address,address,uint256) request) view returns(bytes32)
func (_Auth *AuthCallerSession) HashWithdrawERC20Request(request WithdrawERC20Request) ([32]byte, error) {
	return _Auth.Contract.HashWithdrawERC20Request(&_Auth.CallOpts, request)
}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Auth *AuthCaller) HashWithdrawETHRequest(opts *bind.CallOpts, request WithdrawETHRequest) ([32]byte, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "hashWithdrawETHRequest", request)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Auth *AuthSession) HashWithdrawETHRequest(request WithdrawETHRequest) ([32]byte, error) {
	return _Auth.Contract.HashWithdrawETHRequest(&_Auth.CallOpts, request)
}

// HashWithdrawETHRequest is a free data retrieval call binding the contract method 0x4c24f0b0.
//
// Solidity: function hashWithdrawETHRequest((uint256,address,uint256) request) view returns(bytes32)
func (_Auth *AuthCallerSession) HashWithdrawETHRequest(request WithdrawETHRequest) ([32]byte, error) {
	return _Auth.Contract.HashWithdrawETHRequest(&_Auth.CallOpts, request)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Auth *AuthCaller) Signers(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Auth.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Auth *AuthSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _Auth.Contract.Signers(&_Auth.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers(uint256 ) view returns(address)
func (_Auth *AuthCallerSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _Auth.Contract.Signers(&_Auth.CallOpts, arg0)
}

// AuthRegisterSignersIterator is returned from FilterRegisterSigners and is used to iterate over the raw logs and unpacked data for RegisterSigners events raised by the Auth contract.
type AuthRegisterSignersIterator struct {
	Event *AuthRegisterSigners // Event containing the contract specifics and raw log

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
func (it *AuthRegisterSignersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthRegisterSigners)
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
		it.Event = new(AuthRegisterSigners)
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
func (it *AuthRegisterSignersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuthRegisterSignersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuthRegisterSigners represents a RegisterSigners event raised by the Auth contract.
type AuthRegisterSigners struct {
	Signers []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRegisterSigners is a free log retrieval operation binding the contract event 0x17ef76377f15d668f61d70a53ea5efbe0a4417c081652ede13ec91f1cee880b0.
//
// Solidity: event RegisterSigners(address[] signers)
func (_Auth *AuthFilterer) FilterRegisterSigners(opts *bind.FilterOpts) (*AuthRegisterSignersIterator, error) {

	logs, sub, err := _Auth.contract.FilterLogs(opts, "RegisterSigners")
	if err != nil {
		return nil, err
	}
	return &AuthRegisterSignersIterator{contract: _Auth.contract, event: "RegisterSigners", logs: logs, sub: sub}, nil
}

// WatchRegisterSigners is a free log subscription operation binding the contract event 0x17ef76377f15d668f61d70a53ea5efbe0a4417c081652ede13ec91f1cee880b0.
//
// Solidity: event RegisterSigners(address[] signers)
func (_Auth *AuthFilterer) WatchRegisterSigners(opts *bind.WatchOpts, sink chan<- *AuthRegisterSigners) (event.Subscription, error) {

	logs, sub, err := _Auth.contract.WatchLogs(opts, "RegisterSigners")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuthRegisterSigners)
				if err := _Auth.contract.UnpackLog(event, "RegisterSigners", log); err != nil {
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
func (_Auth *AuthFilterer) ParseRegisterSigners(log types.Log) (*AuthRegisterSigners, error) {
	event := new(AuthRegisterSigners)
	if err := _Auth.contract.UnpackLog(event, "RegisterSigners", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
