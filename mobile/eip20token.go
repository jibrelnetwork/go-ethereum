// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package geth

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// EIP20ABI is the input ABI used to generate the binding from.
const EIP20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialAmount\",\"type\":\"uint256\"},{\"name\":\"_tokenName\",\"type\":\"string\"},{\"name\":\"_decimalUnits\",\"type\":\"uint8\"},{\"name\":\"_tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// EIP20Bin is the compiled bytecode used for deploying new contracts.
const EIP20Bin = `0x6060604052341561000f57600080fd5b6040516107ae3803806107ae833981016040528080519190602001805182019190602001805191906020018051600160a060020a03331660009081526001602052604081208790558690559091019050600383805161007292916020019061009f565b506004805460ff191660ff8416179055600581805161009592916020019061009f565b505050505061013a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100e057805160ff191683800117855561010d565b8280016001018555821561010d579182015b8281111561010d5782518255916020019190600101906100f2565b5061011992915061011d565b5090565b61013791905b808211156101195760008155600101610123565b90565b610665806101496000396000f3006060604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde0381146100b3578063095ea7b31461013d57806318160ddd1461017357806323b872dd1461019857806327e235e3146101c0578063313ce567146101df5780635c6581651461020857806370a082311461022d57806395d89b411461024c578063a9059cbb1461025f578063dd62ed3e14610281575b600080fd5b34156100be57600080fd5b6100c66102a6565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101025780820151838201526020016100ea565b50505050905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014857600080fd5b61015f600160a060020a0360043516602435610344565b604051901515815260200160405180910390f35b341561017e57600080fd5b6101866103b0565b60405190815260200160405180910390f35b34156101a357600080fd5b61015f600160a060020a03600435811690602435166044356103b6565b34156101cb57600080fd5b610186600160a060020a03600435166104bc565b34156101ea57600080fd5b6101f26104ce565b60405160ff909116815260200160405180910390f35b341561021357600080fd5b610186600160a060020a03600435811690602435166104d7565b341561023857600080fd5b610186600160a060020a03600435166104f4565b341561025757600080fd5b6100c661050f565b341561026a57600080fd5b61015f600160a060020a036004351660243561057a565b341561028c57600080fd5b610186600160a060020a036004358116906024351661060e565b60038054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561033c5780601f106103115761010080835404028352916020019161033c565b820191906000526020600020905b81548152906001019060200180831161031f57829003601f168201915b505050505081565b600160a060020a03338116600081815260026020908152604080832094871680845294909152808220859055909291907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a350600192915050565b60005481565b600160a060020a0380841660008181526002602090815260408083203390951683529381528382205492825260019052918220548390108015906103fa5750828110155b151561040557600080fd5b600160a060020a038085166000908152600160205260408082208054870190559187168152208054849003905560001981101561046a57600160a060020a03808616600090815260026020908152604080832033909416835292905220805484900390555b83600160a060020a031685600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8560405190815260200160405180910390a3506001949350505050565b60016020526000908152604090205481565b60045460ff1681565b600260209081526000928352604080842090915290825290205481565b600160a060020a031660009081526001602052604090205490565b60058054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561033c5780601f106103115761010080835404028352916020019161033c565b600160a060020a033316600090815260016020526040812054829010156105a057600080fd5b600160a060020a033381166000818152600160205260408082208054879003905592861680825290839020805486019055917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b600160a060020a039182166000908152600260209081526040808320939094168252919091522054905600a165627a7a7230582029e25a0e3eab2f6ceae984114cbe73a51fa8097006e0ed9c72ff450ac9d489c40029`

// DeployEIP20 deploys a new Ethereum contract, binding an instance of EIP20 to it.
func DeployEIP20(auth *bind.TransactOpts, backend bind.ContractBackend, _initialAmount *big.Int, _tokenName string, _decimalUnits uint8, _tokenSymbol string) (common.Address, *types.Transaction, *EIP20, error) {
	parsed, err := abi.JSON(strings.NewReader(EIP20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EIP20Bin), backend, _initialAmount, _tokenName, _decimalUnits, _tokenSymbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EIP20{EIP20Caller: EIP20Caller{contract: contract}, EIP20Transactor: EIP20Transactor{contract: contract}, EIP20Filterer: EIP20Filterer{contract: contract}}, nil
}

// EIP20 is an auto generated Go binding around an Ethereum contract.
type EIP20 struct {
	EIP20Caller     // Read-only binding to the contract
	EIP20Transactor // Write-only binding to the contract
	EIP20Filterer   // Log filterer for contract events
}

// EIP20Caller is an auto generated read-only Go binding around an Ethereum contract.
type EIP20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP20Session struct {
	Contract     *EIP20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP20CallerSession struct {
	Contract *EIP20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EIP20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP20TransactorSession struct {
	Contract     *EIP20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP20Raw is an auto generated low-level Go binding around an Ethereum contract.
type EIP20Raw struct {
	Contract *EIP20 // Generic contract binding to access the raw methods on
}

// EIP20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP20CallerRaw struct {
	Contract *EIP20Caller // Generic read-only contract binding to access the raw methods on
}

// EIP20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP20TransactorRaw struct {
	Contract *EIP20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP20 creates a new instance of EIP20, bound to a specific deployed contract.
func NewEIP20(address common.Address, backend bind.ContractBackend) (*EIP20, error) {
	contract, err := bindEIP20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP20{EIP20Caller: EIP20Caller{contract: contract}, EIP20Transactor: EIP20Transactor{contract: contract}, EIP20Filterer: EIP20Filterer{contract: contract}}, nil
}

// NewEIP20Caller creates a new read-only instance of EIP20, bound to a specific deployed contract.
func NewEIP20Caller(address common.Address, caller bind.ContractCaller) (*EIP20Caller, error) {
	contract, err := bindEIP20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20Caller{contract: contract}, nil
}

// NewEIP20Transactor creates a new write-only instance of EIP20, bound to a specific deployed contract.
func NewEIP20Transactor(address common.Address, transactor bind.ContractTransactor) (*EIP20Transactor, error) {
	contract, err := bindEIP20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20Transactor{contract: contract}, nil
}

// NewEIP20Filterer creates a new log filterer instance of EIP20, bound to a specific deployed contract.
func NewEIP20Filterer(address common.Address, filterer bind.ContractFilterer) (*EIP20Filterer, error) {
	contract, err := bindEIP20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP20Filterer{contract: contract}, nil
}

// bindEIP20 binds a generic wrapper to an already deployed contract.
func bindEIP20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EIP20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20 *EIP20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EIP20.Contract.EIP20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20 *EIP20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20.Contract.EIP20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20 *EIP20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20.Contract.EIP20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20 *EIP20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EIP20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20 *EIP20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20 *EIP20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20 *EIP20Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20 *EIP20Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowance(&_EIP20.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20 *EIP20CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowance(&_EIP20.CallOpts, _owner, _spender)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed( address,  address) constant returns(uint256)
func (_EIP20 *EIP20Caller) Allowed(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "allowed", arg0, arg1)
	return *ret0, err
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed( address,  address) constant returns(uint256)
func (_EIP20 *EIP20Session) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowed(&_EIP20.CallOpts, arg0, arg1)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed( address,  address) constant returns(uint256)
func (_EIP20 *EIP20CallerSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowed(&_EIP20.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20 *EIP20Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20 *EIP20Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20.Contract.BalanceOf(&_EIP20.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20 *EIP20CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20.Contract.BalanceOf(&_EIP20.CallOpts, _owner)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_EIP20 *EIP20Caller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_EIP20 *EIP20Session) Balances(arg0 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Balances(&_EIP20.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_EIP20 *EIP20CallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Balances(&_EIP20.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EIP20 *EIP20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EIP20 *EIP20Session) Decimals() (uint8, error) {
	return _EIP20.Contract.Decimals(&_EIP20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_EIP20 *EIP20CallerSession) Decimals() (uint8, error) {
	return _EIP20.Contract.Decimals(&_EIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EIP20 *EIP20Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EIP20 *EIP20Session) Name() (string, error) {
	return _EIP20.Contract.Name(&_EIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_EIP20 *EIP20CallerSession) Name() (string, error) {
	return _EIP20.Contract.Name(&_EIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EIP20 *EIP20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EIP20 *EIP20Session) Symbol() (string, error) {
	return _EIP20.Contract.Symbol(&_EIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_EIP20 *EIP20CallerSession) Symbol() (string, error) {
	return _EIP20.Contract.Symbol(&_EIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20 *EIP20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20 *EIP20Session) TotalSupply() (*big.Int, error) {
	return _EIP20.Contract.TotalSupply(&_EIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20 *EIP20CallerSession) TotalSupply() (*big.Int, error) {
	return _EIP20.Contract.TotalSupply(&_EIP20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Approve(&_EIP20.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20 *EIP20TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Approve(&_EIP20.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Transfer(&_EIP20.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Transfer(&_EIP20.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.TransferFrom(&_EIP20.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20 *EIP20TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.TransferFrom(&_EIP20.TransactOpts, _from, _to, _value)
}

// EIP20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EIP20 contract.
type EIP20ApprovalIterator struct {
	Event *EIP20Approval // Event containing the contract specifics and raw log

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
func (it *EIP20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20Approval)
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
		it.Event = new(EIP20Approval)
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
func (it *EIP20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20Approval represents a Approval event raised by the EIP20 contract.
type EIP20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EIP20 *EIP20Filterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*EIP20ApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &EIP20ApprovalIterator{contract: _EIP20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EIP20 *EIP20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EIP20Approval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20Approval)
				if err := _EIP20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// EIP20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EIP20 contract.
type EIP20TransferIterator struct {
	Event *EIP20Transfer // Event containing the contract specifics and raw log

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
func (it *EIP20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20Transfer)
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
		it.Event = new(EIP20Transfer)
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
func (it *EIP20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20Transfer represents a Transfer event raised by the EIP20 contract.
type EIP20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EIP20 *EIP20Filterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*EIP20TransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &EIP20TransferIterator{contract: _EIP20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EIP20 *EIP20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EIP20Transfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20Transfer)
				if err := _EIP20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// EIP20InterfaceABI is the input ABI used to generate the binding from.
const EIP20InterfaceABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// EIP20InterfaceBin is the compiled bytecode used for deploying new contracts.
const EIP20InterfaceBin = `0x`

// DeployEIP20Interface deploys a new Ethereum contract, binding an instance of EIP20Interface to it.
func DeployEIP20Interface(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EIP20Interface, error) {
	parsed, err := abi.JSON(strings.NewReader(EIP20InterfaceABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EIP20InterfaceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EIP20Interface{EIP20InterfaceCaller: EIP20InterfaceCaller{contract: contract}, EIP20InterfaceTransactor: EIP20InterfaceTransactor{contract: contract}, EIP20InterfaceFilterer: EIP20InterfaceFilterer{contract: contract}}, nil
}

// EIP20Interface is an auto generated Go binding around an Ethereum contract.
type EIP20Interface struct {
	EIP20InterfaceCaller     // Read-only binding to the contract
	EIP20InterfaceTransactor // Write-only binding to the contract
	EIP20InterfaceFilterer   // Log filterer for contract events
}

// EIP20InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type EIP20InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP20InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP20InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP20InterfaceSession struct {
	Contract     *EIP20Interface   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP20InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP20InterfaceCallerSession struct {
	Contract *EIP20InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// EIP20InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP20InterfaceTransactorSession struct {
	Contract     *EIP20InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// EIP20InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type EIP20InterfaceRaw struct {
	Contract *EIP20Interface // Generic contract binding to access the raw methods on
}

// EIP20InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP20InterfaceCallerRaw struct {
	Contract *EIP20InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// EIP20InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP20InterfaceTransactorRaw struct {
	Contract *EIP20InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP20Interface creates a new instance of EIP20Interface, bound to a specific deployed contract.
func NewEIP20Interface(address common.Address, backend bind.ContractBackend) (*EIP20Interface, error) {
	contract, err := bindEIP20Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP20Interface{EIP20InterfaceCaller: EIP20InterfaceCaller{contract: contract}, EIP20InterfaceTransactor: EIP20InterfaceTransactor{contract: contract}, EIP20InterfaceFilterer: EIP20InterfaceFilterer{contract: contract}}, nil
}

// NewEIP20InterfaceCaller creates a new read-only instance of EIP20Interface, bound to a specific deployed contract.
func NewEIP20InterfaceCaller(address common.Address, caller bind.ContractCaller) (*EIP20InterfaceCaller, error) {
	contract, err := bindEIP20Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20InterfaceCaller{contract: contract}, nil
}

// NewEIP20InterfaceTransactor creates a new write-only instance of EIP20Interface, bound to a specific deployed contract.
func NewEIP20InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*EIP20InterfaceTransactor, error) {
	contract, err := bindEIP20Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20InterfaceTransactor{contract: contract}, nil
}

// NewEIP20InterfaceFilterer creates a new log filterer instance of EIP20Interface, bound to a specific deployed contract.
func NewEIP20InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*EIP20InterfaceFilterer, error) {
	contract, err := bindEIP20Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP20InterfaceFilterer{contract: contract}, nil
}

// bindEIP20Interface binds a generic wrapper to an already deployed contract.
func bindEIP20Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EIP20InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20Interface *EIP20InterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EIP20Interface.Contract.EIP20InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20Interface *EIP20InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20Interface.Contract.EIP20InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20Interface *EIP20InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20Interface.Contract.EIP20InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20Interface *EIP20InterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EIP20Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20Interface *EIP20InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20Interface *EIP20InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20Interface.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20Interface *EIP20InterfaceCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20Interface.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20Interface *EIP20InterfaceSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20Interface.Contract.Allowance(&_EIP20Interface.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_EIP20Interface *EIP20InterfaceCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20Interface.Contract.Allowance(&_EIP20Interface.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20Interface *EIP20InterfaceCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20Interface.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20Interface *EIP20InterfaceSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20Interface.Contract.BalanceOf(&_EIP20Interface.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_EIP20Interface *EIP20InterfaceCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20Interface.Contract.BalanceOf(&_EIP20Interface.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20Interface *EIP20InterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EIP20Interface.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20Interface *EIP20InterfaceSession) TotalSupply() (*big.Int, error) {
	return _EIP20Interface.Contract.TotalSupply(&_EIP20Interface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_EIP20Interface *EIP20InterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _EIP20Interface.Contract.TotalSupply(&_EIP20Interface.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.Approve(&_EIP20Interface.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.Approve(&_EIP20Interface.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.Transfer(&_EIP20Interface.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.Transfer(&_EIP20Interface.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.TransferFrom(&_EIP20Interface.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_EIP20Interface *EIP20InterfaceTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20Interface.Contract.TransferFrom(&_EIP20Interface.TransactOpts, _from, _to, _value)
}

// EIP20InterfaceApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EIP20Interface contract.
type EIP20InterfaceApprovalIterator struct {
	Event *EIP20InterfaceApproval // Event containing the contract specifics and raw log

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
func (it *EIP20InterfaceApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20InterfaceApproval)
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
		it.Event = new(EIP20InterfaceApproval)
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
func (it *EIP20InterfaceApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20InterfaceApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20InterfaceApproval represents a Approval event raised by the EIP20Interface contract.
type EIP20InterfaceApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EIP20Interface *EIP20InterfaceFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*EIP20InterfaceApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20Interface.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &EIP20InterfaceApprovalIterator{contract: _EIP20Interface.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_EIP20Interface *EIP20InterfaceFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EIP20InterfaceApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20Interface.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20InterfaceApproval)
				if err := _EIP20Interface.contract.UnpackLog(event, "Approval", log); err != nil {
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

// EIP20InterfaceTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EIP20Interface contract.
type EIP20InterfaceTransferIterator struct {
	Event *EIP20InterfaceTransfer // Event containing the contract specifics and raw log

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
func (it *EIP20InterfaceTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20InterfaceTransfer)
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
		it.Event = new(EIP20InterfaceTransfer)
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
func (it *EIP20InterfaceTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20InterfaceTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20InterfaceTransfer represents a Transfer event raised by the EIP20Interface contract.
type EIP20InterfaceTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EIP20Interface *EIP20InterfaceFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*EIP20InterfaceTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20Interface.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &EIP20InterfaceTransferIterator{contract: _EIP20Interface.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(_from indexed address, _to indexed address, _value uint256)
func (_EIP20Interface *EIP20InterfaceFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EIP20InterfaceTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20Interface.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20InterfaceTransfer)
				if err := _EIP20Interface.contract.UnpackLog(event, "Transfer", log); err != nil {
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
