// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package SimpleEscrowEvents

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

// SimpleEscrowEventsMetaData contains all meta data concerning the SimpleEscrowEvents contract.
var SimpleEscrowEventsMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"name\":\"destinationStellarAddress\",\"type\":\"string\"},{\"name\":\"tokenContractAddress\",\"type\":\"string\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"destinationStellarAddress\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"tokenContractAddress\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"Payment\",\"type\":\"event\"}]",
	Sigs: map[string]string{
		"cd485c28": "send(string,string,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b506104e7806100206000396000f3006080604052600436106100405763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663cd485c288114610045575b600080fd5b34801561005157600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100de94369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a99988101979196509182019450925082915084018382808284375094975050933594506100e09350505050565b005b61014682606060405190810160405280602a81526020017f307830303030303030303030303030303030303030303030303030303030303081526020017f303030303030303030300000000000000000000000000000000000000000000081525061032e565b806101b257506101b282606060405190810160405280602a81526020017f307861306238363939316336323138623336633164313964346132653965623081526020017f636533363036656234380000000000000000000000000000000000000000000081525061032e565b151561024557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f746f6b656e436f6e747261637441646472657373206e6f7420737570706f727460448201527f6564000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b816040518082805190602001908083835b602083106102755780518252601f199092019160209182019101610256565b51815160209384036101000a6000190180199092169116179052604051919093018190038120885190955088945090928392508401908083835b602083106102ce5780518252601f1990920191602091820191016102af565b51815160209384036101000a60001901801990921691161790526040805192909401829003822088835293519395507fa65615f26a7f49d1df9b25d7019bf788d420564dbe08ad9c06d70ef5519c05ff94509083900301919050a3505050565b6000816040516020018082805190602001908083835b602083106103635780518252601f199092019160209182019101610344565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b602083106103c65780518252601f1990920191602091820191016103a7565b51815160209384036101000a60001901801990921691161790526040519190930181900381208851909550889450908301928392508401908083835b602083106104215780518252601f199092019160209182019101610402565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b602083106104845780518252601f199092019160209182019101610465565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209390931496955050505050505600a165627a7a72305820082bf9fb96f0a3d77f788abb0b773db6933df158793648aa1cf0ac2de3ffe51d0029",
}

// SimpleEscrowEventsABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleEscrowEventsMetaData.ABI instead.
var SimpleEscrowEventsABI = SimpleEscrowEventsMetaData.ABI

// Deprecated: Use SimpleEscrowEventsMetaData.Sigs instead.
// SimpleEscrowEventsFuncSigs maps the 4-byte function signature to its string representation.
var SimpleEscrowEventsFuncSigs = SimpleEscrowEventsMetaData.Sigs

// SimpleEscrowEventsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimpleEscrowEventsMetaData.Bin instead.
var SimpleEscrowEventsBin = SimpleEscrowEventsMetaData.Bin

// DeploySimpleEscrowEvents deploys a new Ethereum contract, binding an instance of SimpleEscrowEvents to it.
func DeploySimpleEscrowEvents(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleEscrowEvents, error) {
	parsed, err := SimpleEscrowEventsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleEscrowEventsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleEscrowEvents{SimpleEscrowEventsCaller: SimpleEscrowEventsCaller{contract: contract}, SimpleEscrowEventsTransactor: SimpleEscrowEventsTransactor{contract: contract}, SimpleEscrowEventsFilterer: SimpleEscrowEventsFilterer{contract: contract}}, nil
}

// SimpleEscrowEvents is an auto generated Go binding around an Ethereum contract.
type SimpleEscrowEvents struct {
	SimpleEscrowEventsCaller     // Read-only binding to the contract
	SimpleEscrowEventsTransactor // Write-only binding to the contract
	SimpleEscrowEventsFilterer   // Log filterer for contract events
}

// SimpleEscrowEventsCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleEscrowEventsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleEscrowEventsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleEscrowEventsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleEscrowEventsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleEscrowEventsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleEscrowEventsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleEscrowEventsSession struct {
	Contract     *SimpleEscrowEvents // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SimpleEscrowEventsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleEscrowEventsCallerSession struct {
	Contract *SimpleEscrowEventsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SimpleEscrowEventsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleEscrowEventsTransactorSession struct {
	Contract     *SimpleEscrowEventsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SimpleEscrowEventsRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleEscrowEventsRaw struct {
	Contract *SimpleEscrowEvents // Generic contract binding to access the raw methods on
}

// SimpleEscrowEventsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleEscrowEventsCallerRaw struct {
	Contract *SimpleEscrowEventsCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleEscrowEventsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleEscrowEventsTransactorRaw struct {
	Contract *SimpleEscrowEventsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleEscrowEvents creates a new instance of SimpleEscrowEvents, bound to a specific deployed contract.
func NewSimpleEscrowEvents(address common.Address, backend bind.ContractBackend) (*SimpleEscrowEvents, error) {
	contract, err := bindSimpleEscrowEvents(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEvents{SimpleEscrowEventsCaller: SimpleEscrowEventsCaller{contract: contract}, SimpleEscrowEventsTransactor: SimpleEscrowEventsTransactor{contract: contract}, SimpleEscrowEventsFilterer: SimpleEscrowEventsFilterer{contract: contract}}, nil
}

// NewSimpleEscrowEventsCaller creates a new read-only instance of SimpleEscrowEvents, bound to a specific deployed contract.
func NewSimpleEscrowEventsCaller(address common.Address, caller bind.ContractCaller) (*SimpleEscrowEventsCaller, error) {
	contract, err := bindSimpleEscrowEvents(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEventsCaller{contract: contract}, nil
}

// NewSimpleEscrowEventsTransactor creates a new write-only instance of SimpleEscrowEvents, bound to a specific deployed contract.
func NewSimpleEscrowEventsTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleEscrowEventsTransactor, error) {
	contract, err := bindSimpleEscrowEvents(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEventsTransactor{contract: contract}, nil
}

// NewSimpleEscrowEventsFilterer creates a new log filterer instance of SimpleEscrowEvents, bound to a specific deployed contract.
func NewSimpleEscrowEventsFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleEscrowEventsFilterer, error) {
	contract, err := bindSimpleEscrowEvents(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEventsFilterer{contract: contract}, nil
}

// bindSimpleEscrowEvents binds a generic wrapper to an already deployed contract.
func bindSimpleEscrowEvents(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleEscrowEventsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleEscrowEvents *SimpleEscrowEventsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleEscrowEvents.Contract.SimpleEscrowEventsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleEscrowEvents *SimpleEscrowEventsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.SimpleEscrowEventsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleEscrowEvents *SimpleEscrowEventsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.SimpleEscrowEventsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleEscrowEvents *SimpleEscrowEventsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleEscrowEvents.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.contract.Transact(opts, method, params...)
}

// Send is a paid mutator transaction binding the contract method 0xcd485c28.
//
// Solidity: function send(string destinationStellarAddress, string tokenContractAddress, uint256 tokenAmount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactor) Send(opts *bind.TransactOpts, destinationStellarAddress string, tokenContractAddress string, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.contract.Transact(opts, "send", destinationStellarAddress, tokenContractAddress, tokenAmount)
}

// Send is a paid mutator transaction binding the contract method 0xcd485c28.
//
// Solidity: function send(string destinationStellarAddress, string tokenContractAddress, uint256 tokenAmount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsSession) Send(destinationStellarAddress string, tokenContractAddress string, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.Send(&_SimpleEscrowEvents.TransactOpts, destinationStellarAddress, tokenContractAddress, tokenAmount)
}

// Send is a paid mutator transaction binding the contract method 0xcd485c28.
//
// Solidity: function send(string destinationStellarAddress, string tokenContractAddress, uint256 tokenAmount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactorSession) Send(destinationStellarAddress string, tokenContractAddress string, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.Send(&_SimpleEscrowEvents.TransactOpts, destinationStellarAddress, tokenContractAddress, tokenAmount)
}

// SimpleEscrowEventsPaymentIterator is returned from FilterPayment and is used to iterate over the raw logs and unpacked data for Payment events raised by the SimpleEscrowEvents contract.
type SimpleEscrowEventsPaymentIterator struct {
	Event *SimpleEscrowEventsPayment // Event containing the contract specifics and raw log

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
func (it *SimpleEscrowEventsPaymentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleEscrowEventsPayment)
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
		it.Event = new(SimpleEscrowEventsPayment)
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
func (it *SimpleEscrowEventsPaymentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleEscrowEventsPaymentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleEscrowEventsPayment represents a Payment event raised by the SimpleEscrowEvents contract.
type SimpleEscrowEventsPayment struct {
	DestinationStellarAddress common.Hash
	TokenContractAddress      common.Hash
	TokenAmount               *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterPayment is a free log retrieval operation binding the contract event 0xa65615f26a7f49d1df9b25d7019bf788d420564dbe08ad9c06d70ef5519c05ff.
//
// Solidity: event Payment(string indexed destinationStellarAddress, string indexed tokenContractAddress, uint256 tokenAmount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) FilterPayment(opts *bind.FilterOpts, destinationStellarAddress []string, tokenContractAddress []string) (*SimpleEscrowEventsPaymentIterator, error) {

	var destinationStellarAddressRule []interface{}
	for _, destinationStellarAddressItem := range destinationStellarAddress {
		destinationStellarAddressRule = append(destinationStellarAddressRule, destinationStellarAddressItem)
	}
	var tokenContractAddressRule []interface{}
	for _, tokenContractAddressItem := range tokenContractAddress {
		tokenContractAddressRule = append(tokenContractAddressRule, tokenContractAddressItem)
	}

	logs, sub, err := _SimpleEscrowEvents.contract.FilterLogs(opts, "Payment", destinationStellarAddressRule, tokenContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEventsPaymentIterator{contract: _SimpleEscrowEvents.contract, event: "Payment", logs: logs, sub: sub}, nil
}

// WatchPayment is a free log subscription operation binding the contract event 0xa65615f26a7f49d1df9b25d7019bf788d420564dbe08ad9c06d70ef5519c05ff.
//
// Solidity: event Payment(string indexed destinationStellarAddress, string indexed tokenContractAddress, uint256 tokenAmount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) WatchPayment(opts *bind.WatchOpts, sink chan<- *SimpleEscrowEventsPayment, destinationStellarAddress []string, tokenContractAddress []string) (event.Subscription, error) {

	var destinationStellarAddressRule []interface{}
	for _, destinationStellarAddressItem := range destinationStellarAddress {
		destinationStellarAddressRule = append(destinationStellarAddressRule, destinationStellarAddressItem)
	}
	var tokenContractAddressRule []interface{}
	for _, tokenContractAddressItem := range tokenContractAddress {
		tokenContractAddressRule = append(tokenContractAddressRule, tokenContractAddressItem)
	}

	logs, sub, err := _SimpleEscrowEvents.contract.WatchLogs(opts, "Payment", destinationStellarAddressRule, tokenContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleEscrowEventsPayment)
				if err := _SimpleEscrowEvents.contract.UnpackLog(event, "Payment", log); err != nil {
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

// ParsePayment is a log parse operation binding the contract event 0xa65615f26a7f49d1df9b25d7019bf788d420564dbe08ad9c06d70ef5519c05ff.
//
// Solidity: event Payment(string indexed destinationStellarAddress, string indexed tokenContractAddress, uint256 tokenAmount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) ParsePayment(log types.Log) (*SimpleEscrowEventsPayment, error) {
	event := new(SimpleEscrowEventsPayment)
	if err := _SimpleEscrowEvents.contract.UnpackLog(event, "Payment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
