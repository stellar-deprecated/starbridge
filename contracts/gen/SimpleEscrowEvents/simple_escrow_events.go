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
	ABI: "[{\"constant\":false,\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"a\",\"type\":\"string\"},{\"name\":\"b\",\"type\":\"string\"}],\"name\":\"isStringsEqual\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractAddress\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Payment\",\"type\":\"event\"}]",
	Sigs: map[string]string{
		"ede429bc": "isStringsEqual(string,string)",
		"ccc95e03": "send(string,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b506104df806100206000396000f30060806040526004361061004b5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663ccc95e038114610050578063ede429bc146100ad575b600080fd5b34801561005c57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100ab94369492936024939284019190819084018382808284375094975050933594506101589350505050565b005b3480156100b957600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261014494369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506103269650505050505050565b604080519115158252519081900360200190f35b6101be82606060405190810160405280602a81526020017f307830303030303030303030303030303030303030303030303030303030303081526020017f3030303030303030303000000000000000000000000000000000000000000000815250610326565b8061022a575061022a82606060405190810160405280602a81526020017f307861306238363939316336323138623336633164313964346132653965623081526020017f6365333630366562343800000000000000000000000000000000000000000000815250610326565b151561029757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f636f6e747261637441646472657373206e6f7420737570706f72746564000000604482015290519081900360640190fd5b816040518082805190602001908083835b602083106102c75780518252601f1990920191602091820191016102a8565b51815160209384036101000a60001901801990921691161790526040805192909401829003822087835293519395507f10bc7b26b8e5a9e4746a8ab34931348ab6bb1caeeb7f7d4aec1febe619194b0b94509083900301919050a25050565b6000816040516020018082805190602001908083835b6020831061035b5780518252601f19909201916020918201910161033c565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b602083106103be5780518252601f19909201916020918201910161039f565b51815160209384036101000a60001901801990921691161790526040519190930181900381208851909550889450908301928392508401908083835b602083106104195780518252601f1990920191602091820191016103fa565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b6020831061047c5780518252601f19909201916020918201910161045d565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209390931496955050505050505600a165627a7a72305820c45564e248c77dcd6bfe18079c139f4fc15842cf374b7511ff8046e76a3893160029",
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

// IsStringsEqual is a free data retrieval call binding the contract method 0xede429bc.
//
// Solidity: function isStringsEqual(string a, string b) pure returns(bool)
func (_SimpleEscrowEvents *SimpleEscrowEventsCaller) IsStringsEqual(opts *bind.CallOpts, a string, b string) (bool, error) {
	var out []interface{}
	err := _SimpleEscrowEvents.contract.Call(opts, &out, "isStringsEqual", a, b)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStringsEqual is a free data retrieval call binding the contract method 0xede429bc.
//
// Solidity: function isStringsEqual(string a, string b) pure returns(bool)
func (_SimpleEscrowEvents *SimpleEscrowEventsSession) IsStringsEqual(a string, b string) (bool, error) {
	return _SimpleEscrowEvents.Contract.IsStringsEqual(&_SimpleEscrowEvents.CallOpts, a, b)
}

// IsStringsEqual is a free data retrieval call binding the contract method 0xede429bc.
//
// Solidity: function isStringsEqual(string a, string b) pure returns(bool)
func (_SimpleEscrowEvents *SimpleEscrowEventsCallerSession) IsStringsEqual(a string, b string) (bool, error) {
	return _SimpleEscrowEvents.Contract.IsStringsEqual(&_SimpleEscrowEvents.CallOpts, a, b)
}

// Send is a paid mutator transaction binding the contract method 0xccc95e03.
//
// Solidity: function send(string contractAddress, uint256 amount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactor) Send(opts *bind.TransactOpts, contractAddress string, amount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.contract.Transact(opts, "send", contractAddress, amount)
}

// Send is a paid mutator transaction binding the contract method 0xccc95e03.
//
// Solidity: function send(string contractAddress, uint256 amount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsSession) Send(contractAddress string, amount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.Send(&_SimpleEscrowEvents.TransactOpts, contractAddress, amount)
}

// Send is a paid mutator transaction binding the contract method 0xccc95e03.
//
// Solidity: function send(string contractAddress, uint256 amount) returns()
func (_SimpleEscrowEvents *SimpleEscrowEventsTransactorSession) Send(contractAddress string, amount *big.Int) (*types.Transaction, error) {
	return _SimpleEscrowEvents.Contract.Send(&_SimpleEscrowEvents.TransactOpts, contractAddress, amount)
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
	ContractAddress common.Hash
	Amount          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPayment is a free log retrieval operation binding the contract event 0x10bc7b26b8e5a9e4746a8ab34931348ab6bb1caeeb7f7d4aec1febe619194b0b.
//
// Solidity: event Payment(string indexed contractAddress, uint256 amount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) FilterPayment(opts *bind.FilterOpts, contractAddress []string) (*SimpleEscrowEventsPaymentIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _SimpleEscrowEvents.contract.FilterLogs(opts, "Payment", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &SimpleEscrowEventsPaymentIterator{contract: _SimpleEscrowEvents.contract, event: "Payment", logs: logs, sub: sub}, nil
}

// WatchPayment is a free log subscription operation binding the contract event 0x10bc7b26b8e5a9e4746a8ab34931348ab6bb1caeeb7f7d4aec1febe619194b0b.
//
// Solidity: event Payment(string indexed contractAddress, uint256 amount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) WatchPayment(opts *bind.WatchOpts, sink chan<- *SimpleEscrowEventsPayment, contractAddress []string) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _SimpleEscrowEvents.contract.WatchLogs(opts, "Payment", contractAddressRule)
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

// ParsePayment is a log parse operation binding the contract event 0x10bc7b26b8e5a9e4746a8ab34931348ab6bb1caeeb7f7d4aec1febe619194b0b.
//
// Solidity: event Payment(string indexed contractAddress, uint256 amount)
func (_SimpleEscrowEvents *SimpleEscrowEventsFilterer) ParsePayment(log types.Log) (*SimpleEscrowEventsPayment, error) {
	event := new(SimpleEscrowEventsPayment)
	if err := _SimpleEscrowEvents.contract.UnpackLog(event, "Payment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
