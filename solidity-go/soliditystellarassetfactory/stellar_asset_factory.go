// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package soliditystellarassetfactory

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

// StellarAssetFactoryMetaData contains all meta data concerning the StellarAssetFactory contract.
var StellarAssetFactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"CreateStellarAsset\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isStellarAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nameHashToAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// StellarAssetFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use StellarAssetFactoryMetaData.ABI instead.
var StellarAssetFactoryABI = StellarAssetFactoryMetaData.ABI

// StellarAssetFactory is an auto generated Go binding around an Ethereum contract.
type StellarAssetFactory struct {
	StellarAssetFactoryCaller     // Read-only binding to the contract
	StellarAssetFactoryTransactor // Write-only binding to the contract
	StellarAssetFactoryFilterer   // Log filterer for contract events
}

// StellarAssetFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type StellarAssetFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StellarAssetFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StellarAssetFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StellarAssetFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StellarAssetFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StellarAssetFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StellarAssetFactorySession struct {
	Contract     *StellarAssetFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// StellarAssetFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StellarAssetFactoryCallerSession struct {
	Contract *StellarAssetFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// StellarAssetFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StellarAssetFactoryTransactorSession struct {
	Contract     *StellarAssetFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// StellarAssetFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type StellarAssetFactoryRaw struct {
	Contract *StellarAssetFactory // Generic contract binding to access the raw methods on
}

// StellarAssetFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StellarAssetFactoryCallerRaw struct {
	Contract *StellarAssetFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// StellarAssetFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StellarAssetFactoryTransactorRaw struct {
	Contract *StellarAssetFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStellarAssetFactory creates a new instance of StellarAssetFactory, bound to a specific deployed contract.
func NewStellarAssetFactory(address common.Address, backend bind.ContractBackend) (*StellarAssetFactory, error) {
	contract, err := bindStellarAssetFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StellarAssetFactory{StellarAssetFactoryCaller: StellarAssetFactoryCaller{contract: contract}, StellarAssetFactoryTransactor: StellarAssetFactoryTransactor{contract: contract}, StellarAssetFactoryFilterer: StellarAssetFactoryFilterer{contract: contract}}, nil
}

// NewStellarAssetFactoryCaller creates a new read-only instance of StellarAssetFactory, bound to a specific deployed contract.
func NewStellarAssetFactoryCaller(address common.Address, caller bind.ContractCaller) (*StellarAssetFactoryCaller, error) {
	contract, err := bindStellarAssetFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StellarAssetFactoryCaller{contract: contract}, nil
}

// NewStellarAssetFactoryTransactor creates a new write-only instance of StellarAssetFactory, bound to a specific deployed contract.
func NewStellarAssetFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*StellarAssetFactoryTransactor, error) {
	contract, err := bindStellarAssetFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StellarAssetFactoryTransactor{contract: contract}, nil
}

// NewStellarAssetFactoryFilterer creates a new log filterer instance of StellarAssetFactory, bound to a specific deployed contract.
func NewStellarAssetFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*StellarAssetFactoryFilterer, error) {
	contract, err := bindStellarAssetFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StellarAssetFactoryFilterer{contract: contract}, nil
}

// bindStellarAssetFactory binds a generic wrapper to an already deployed contract.
func bindStellarAssetFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StellarAssetFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StellarAssetFactory *StellarAssetFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StellarAssetFactory.Contract.StellarAssetFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StellarAssetFactory *StellarAssetFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StellarAssetFactory.Contract.StellarAssetFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StellarAssetFactory *StellarAssetFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StellarAssetFactory.Contract.StellarAssetFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StellarAssetFactory *StellarAssetFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StellarAssetFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StellarAssetFactory *StellarAssetFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StellarAssetFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StellarAssetFactory *StellarAssetFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StellarAssetFactory.Contract.contract.Transact(opts, method, params...)
}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_StellarAssetFactory *StellarAssetFactoryCaller) IsStellarAsset(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StellarAssetFactory.contract.Call(opts, &out, "isStellarAsset", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_StellarAssetFactory *StellarAssetFactorySession) IsStellarAsset(arg0 common.Address) (bool, error) {
	return _StellarAssetFactory.Contract.IsStellarAsset(&_StellarAssetFactory.CallOpts, arg0)
}

// IsStellarAsset is a free data retrieval call binding the contract method 0x453c6d97.
//
// Solidity: function isStellarAsset(address ) view returns(bool)
func (_StellarAssetFactory *StellarAssetFactoryCallerSession) IsStellarAsset(arg0 common.Address) (bool, error) {
	return _StellarAssetFactory.Contract.IsStellarAsset(&_StellarAssetFactory.CallOpts, arg0)
}

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_StellarAssetFactory *StellarAssetFactoryCaller) NameHashToAsset(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _StellarAssetFactory.contract.Call(opts, &out, "nameHashToAsset", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_StellarAssetFactory *StellarAssetFactorySession) NameHashToAsset(arg0 [32]byte) (common.Address, error) {
	return _StellarAssetFactory.Contract.NameHashToAsset(&_StellarAssetFactory.CallOpts, arg0)
}

// NameHashToAsset is a free data retrieval call binding the contract method 0xf278cbc4.
//
// Solidity: function nameHashToAsset(bytes32 ) view returns(address)
func (_StellarAssetFactory *StellarAssetFactoryCallerSession) NameHashToAsset(arg0 [32]byte) (common.Address, error) {
	return _StellarAssetFactory.Contract.NameHashToAsset(&_StellarAssetFactory.CallOpts, arg0)
}

// StellarAssetFactoryCreateStellarAssetIterator is returned from FilterCreateStellarAsset and is used to iterate over the raw logs and unpacked data for CreateStellarAsset events raised by the StellarAssetFactory contract.
type StellarAssetFactoryCreateStellarAssetIterator struct {
	Event *StellarAssetFactoryCreateStellarAsset // Event containing the contract specifics and raw log

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
func (it *StellarAssetFactoryCreateStellarAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StellarAssetFactoryCreateStellarAsset)
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
		it.Event = new(StellarAssetFactoryCreateStellarAsset)
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
func (it *StellarAssetFactoryCreateStellarAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StellarAssetFactoryCreateStellarAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StellarAssetFactoryCreateStellarAsset represents a CreateStellarAsset event raised by the StellarAssetFactory contract.
type StellarAssetFactoryCreateStellarAsset struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCreateStellarAsset is a free log retrieval operation binding the contract event 0x103935973961d03366259eb4fc092bbc84b0812cba6d190f70d03f5f7bbf4cc2.
//
// Solidity: event CreateStellarAsset(address asset)
func (_StellarAssetFactory *StellarAssetFactoryFilterer) FilterCreateStellarAsset(opts *bind.FilterOpts) (*StellarAssetFactoryCreateStellarAssetIterator, error) {

	logs, sub, err := _StellarAssetFactory.contract.FilterLogs(opts, "CreateStellarAsset")
	if err != nil {
		return nil, err
	}
	return &StellarAssetFactoryCreateStellarAssetIterator{contract: _StellarAssetFactory.contract, event: "CreateStellarAsset", logs: logs, sub: sub}, nil
}

// WatchCreateStellarAsset is a free log subscription operation binding the contract event 0x103935973961d03366259eb4fc092bbc84b0812cba6d190f70d03f5f7bbf4cc2.
//
// Solidity: event CreateStellarAsset(address asset)
func (_StellarAssetFactory *StellarAssetFactoryFilterer) WatchCreateStellarAsset(opts *bind.WatchOpts, sink chan<- *StellarAssetFactoryCreateStellarAsset) (event.Subscription, error) {

	logs, sub, err := _StellarAssetFactory.contract.WatchLogs(opts, "CreateStellarAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StellarAssetFactoryCreateStellarAsset)
				if err := _StellarAssetFactory.contract.UnpackLog(event, "CreateStellarAsset", log); err != nil {
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
func (_StellarAssetFactory *StellarAssetFactoryFilterer) ParseCreateStellarAsset(log types.Log) (*StellarAssetFactoryCreateStellarAsset, error) {
	event := new(StellarAssetFactoryCreateStellarAsset)
	if err := _StellarAssetFactory.contract.UnpackLog(event, "CreateStellarAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
