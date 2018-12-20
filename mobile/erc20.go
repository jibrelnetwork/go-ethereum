package geth

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jibrel.network/EIP20/contracts"
)

type Contract struct {
	contract *contracts.EIP20
	CallOpts *CallOpts
	Address  *Address
}

type TransactionSigner interface {
	Sign(contract *Contract, address *Address, tx *Transaction) *Transaction
}

func NewContract(address *Address, client *EthereumClient) (contract *Contract, _ error) {
	eip20, err := contracts.NewEIP20(address.address, client.client)
	if err != nil {
		return nil, err
	}
	return &Contract{
		contract: eip20,
		CallOpts: &CallOpts{
			opts: bind.CallOpts{Pending: false},
		},
		Address: address,
	}, nil
}

func (c *Contract) Name(opts *CallOpts) (string, error) {
	return c.contract.Name(c.getCallOpts(opts))
}

func (c *Contract) Symbol(opts *CallOpts) (string, error) {
	return c.contract.Symbol(c.getCallOpts(opts))
}

func (c *Contract) Decimals(opts *CallOpts) (*BigInt, error) {
	result, err := c.contract.Decimals(c.getCallOpts(opts))
	if err != nil {
		return nil, err
	}

	return &BigInt{big.NewInt(int64(result))}, nil
}

func (c *Contract) getCallOpts(opts *CallOpts) (result *bind.CallOpts) {
	if opts == nil {
		return &(c.CallOpts.opts)
	}
	return &opts.opts
}

// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (c *Contract) Allowance(owner *Address, spender *Address, opts *CallOpts) (*BigInt, error) {
	value, err := c.contract.Allowance(c.getCallOpts(opts), owner.address, spender.address)
	if err != nil {
		return nil, err
	}
	return &BigInt{value}, nil
}

// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (c *Contract) BalanceOf(owner *Address, opts *CallOpts) (*BigInt, error) {
	value, err := c.contract.BalanceOf(c.getCallOpts(opts), owner.address)
	if err != nil {
		return nil, err
	}
	return &BigInt{value}, nil
}

// Solidity: function totalSupply() constant returns(uint256)
func (c *Contract) TotalSupply(opts *CallOpts) (*BigInt, error) {
	value, err := c.contract.TotalSupply(c.getCallOpts(opts))
	if err != nil {
		return nil, err
	}
	return &BigInt{value}, nil
}

// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (c *Contract) Approve(opts *TransactOpts, spender *Address, value *BigInt) (*Transaction, error) {
	tx, err := c.contract.Approve(&opts.opts, spender.address, value.bigint)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (c *Contract) Transfer(opts *TransactOpts, to *Address, value *BigInt) (*Transaction, error) {
	tx, err := c.contract.Transfer(&opts.opts, to.address, value.bigint)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx}, nil
}

// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (c *Contract) TransferFrom(opts *TransactOpts, from *Address, to *Address, value *BigInt) (*Transaction, error) {
	tx, err := c.contract.TransferFrom(&opts.opts, from.address, to.address, value.bigint)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx}, nil
}

func (c *Contract) EstimateGasPrice(ctx *Context, opts *TransactOpts) {
	//c.contract.ERC20Caller.contract
}

// Context and signer should not be nil, other params may be nil
func NewTransactOpts(ctx *Context, signer Signer, from *Address, nonce *BigInt, value *BigInt, gasPrice *BigInt, gasLimit *BigInt) *TransactOpts {
	opts := bind.TransactOpts{Context: ctx.context}
	if from != nil {
		opts.From = from.address
	}
	if nonce != nil {
		opts.Nonce = nonce.bigint
	}
	if value != nil {
		opts.Value = value.bigint
	}
	if gasPrice != nil {
		opts.GasPrice = gasPrice.bigint
	}
	if gasLimit != nil {
		opts.GasLimit = uint64(gasLimit.GetInt64())
	}
	result := TransactOpts{opts}
	result.SetSigner(signer)
	return &result
}

type ContractABI struct {
	abi *abi.ABI
}

func NewContractABI(abiString string) (*ContractABI, error) {
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	return &ContractABI{abi: &abi}, nil
}

func NewEIP20ABI() (*ContractABI, error) {
	return NewContractABI(contracts.EIP20ABI)
}

func (abi *ContractABI) PackArguments(method string, params *Parameters) ([]byte, error) {
	return abi.abi.Pack(method, params.params...)
}

func (abi *ContractABI) UnpackString(method string, output []byte) (string, error) {
	var result = new(string)
	err := abi.abi.Unpack(result, method, output)
	return *result, err
}

func (abi *ContractABI) UnpackUInt8(method string, output []byte) (uint8, error) {
	var result = new(uint8)
	err := abi.abi.Unpack(result, method, output)
	return *result, err
}

func (abi *ContractABI) UnpackTransferFrom(topics *Hashes, data []byte) (*TransferEvent, error) {
	var event = new(contracts.ERC20Transfer)
	err := abi.Unpack(event, "Transfer", topics.hashes, data)
	if err != nil {
		return nil, err
	}
	return &TransferEvent{event}, nil
}

func (abi *ContractABI) UnpackTransfer(logWrapper *Log) (*TransferEvent, error) {
	var event = new(contracts.ERC20Transfer)
	err := abi.UnpackLog(event, "Transfer", *logWrapper.log)
	if err != nil {
		return nil, err
	}
	event.Raw = *logWrapper.log
	return &TransferEvent{event}, nil
}

func (abi *ContractABI) UnpackLog(out interface{}, event string, log types.Log) error {
	return abi.Unpack(out, event, log.Topics, log.Data)
}

func (cAbi *ContractABI) Unpack(out interface{}, event string, topics []common.Hash, data []byte) error {
	if len(data) > 0 {
		if err := cAbi.abi.Unpack(out, event, data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range cAbi.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return bind.ParseTopics(out, indexed, topics[1:])
}

type Parameters struct {
	params []interface{}
}

func NewParameters(size int) *Parameters {
	p := make([]interface{}, size)
	return &Parameters{p}
}

func (p *Parameters) get(index int) (interface{}, error) {
	if index < 0 || index >= len(p.params) {
		return nil, errors.New("index out of bounds")
	}
	return p.params[index], nil
}

func (p *Parameters) set(index int, param interface{}) error {
	if index < 0 || index >= len(p.params) {
		return errors.New("index out of bounds")
	}
	p.params[index] = param
	return nil
}

func (p *Parameters) GetString(index int) (string, error) {
	item, err := p.get(index)
	if err != nil {
		return "", err
	}
	return item.(string), err
}

func (p *Parameters) SetString(index int, str string) error {
	return p.set(index, str)
}

func (p *Parameters) GetData(index int) ([]byte, error) {
	item, err := p.get(index)
	if err != nil {
		return nil, err
	}
	return item.([]byte), err
}

func (p *Parameters) SetData(index int, data []byte) error {
	return p.set(index, data)
}

func (p *Parameters) GetBigInt(index int) (*BigInt, error) {
	item, err := p.get(index)
	if err != nil {
		return nil, err
	}
	bigInt := item.(big.Int)
	return &BigInt{&bigInt}, err
}

func (p *Parameters) SetBigInt(index int, bigInt *BigInt) error {
	return p.set(index, bigInt.bigint)
}

func (p *Parameters) GetAddress(index int) (*Address, error) {
	item, err := p.get(index)
	if err != nil {
		return nil, err
	}
	return &Address{item.(common.Address)}, err
}

func (p *Parameters) SetAddress(index int, address *Address) error {
	return p.set(index, address.address)
}

type TransferEvent struct {
	event *contracts.ERC20Transfer
}

func (t *TransferEvent) FromAddress() *Address {
	return &Address{t.event.From}
}

func (t *TransferEvent) ToAddress() *Address {
	return &Address{t.event.To}
}

func (t *TransferEvent) Value() *BigInt {
	return &BigInt{t.event.Value}
}

func (t *TransferEvent) Logs() *Log {
	return &Log{&t.event.Raw}
}
