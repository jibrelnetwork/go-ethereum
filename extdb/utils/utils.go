package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/token"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

type BlockChain interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header

	Config() *params.ChainConfig

	StateAt(root common.Hash) (*state.StateDB, error)
}

type ExtDB interface {
	WriteTokenBalance(tokenBalance *exttypes.TokenBalance) error
}

var (
	erc20_token_abi_str            string
	erc20_approval_event_signature []byte
	erc20_transfer_event_signature []byte
	mint_event_signature           []byte
	burn_event_signature           []byte
	jntmint_event_signature        []byte
	jntburn_event_signature        []byte
)

func init() {
	// ERC-20 ABI
	erc20_token_abi_str = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`
	// Transfer(address,address,uint256)
	erc20_transfer_event_signature, _ = hex.DecodeString("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	// Approval(address,address,uint256)
	erc20_approval_event_signature, _ = hex.DecodeString("8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
	// Mint(address,uint256)
	mint_event_signature, _ = hex.DecodeString("0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885")
	// Burn(address,uint256)
	burn_event_signature, _ = hex.DecodeString("cc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5")
	// MintEvent(address,uint256)
	jntmint_event_signature, _ = hex.DecodeString("3fffaa5804a26fcec0d70b1d0fb0a2d0031df3a5f9c8af2127c2f4360e97b463")
	// BurnEvent(address,uint256)
	jntburn_event_signature, _ = hex.DecodeString("512586160ebd4dc6945ba9ec5d21a1f723f26f3c7aa36cdffb6818d4e7b88030")

	log.Info("extdb::utils::init()")
}

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}

func get_ERC20_token_balance_from_EVM(bc BlockChain, statedb *state.StateDB, block *types.Block, contract_address, querying_addr *common.Address) (*big.Int, error) {
	var (
		err             error
		input           []byte
		erc20_token_abi abi.ABI
	)

	vm_config := vm.Config{
		Debug: false,
	}

	erc20_token_abi, err = abi.JSON(strings.NewReader(erc20_token_abi_str))
	bc_config := bc.Config()
	value := big.NewInt(0)
	gas_limit := uint64(50000000)
	gas_price := big.NewInt(1)
	fake_src_addr := common.HexToAddress("8999999999999999999999999999999999999998")
	fake_balance := big.NewInt(0)
	fake_balance.SetString("9999999999999999999999999999", 10)

	// Add a fake account with huge balance so we have money to pay for gas to execute instructions on the EVM
	statedb.AddBalance(fake_src_addr, fake_balance)

	// Encode input for retrieving token balance
	input, err = erc20_token_abi.Pack("balanceOf", querying_addr)
	if err != nil {
		fatalf("Can't pack balanceOf input: %v", err)
	}

	// Getting token holder balance
	msg := types.NewMessage(fake_src_addr, contract_address, 0, value, gas_limit, gas_price, input, false)
	evm_context := token.NewEVMContext(msg, block.Header(), bc, nil)
	evm := vm.NewEVM(evm_context, statedb, bc_config, vm_config)
	gp := new(token.GasPool).AddGas(math.MaxUint64)

	ret, gas, failed, err := token.ApplyMessage(evm, msg, gp)
	if failed {
		log.Debug(fmt.Sprintf("get_ERC20_token_balance: vm err for symbol: %v, failed=%v", err, failed))
		return nil, fmt.Errorf("vm err")
	}

	if err != nil {
		log.Debug(fmt.Sprintf("get_ERC20_token_balance: getting 'balanceOf' caused error in vm: %v", err))
		return nil, err
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("len(ret)==0")
	}

	balance := big.NewInt(0)
	if !((err != nil) || (failed)) {
		err = erc20_token_abi.Unpack(&balance, "balanceOf", ret)
		if err != nil {
			fatalf("Can't upack balanceOf output from the EVM: %v", err)
		}
	}

	_ = gas

	return balance, nil
}

func get_ERC20_token_decimals_from_EVM(bc BlockChain, statedb *state.StateDB, block *types.Block, contract_address *common.Address) (*uint8, error) {
	var (
		err             error
		input           []byte
		erc20_token_abi abi.ABI
	)

	vm_config := vm.Config{
		Debug: false,
	}

	erc20_token_abi, err = abi.JSON(strings.NewReader(erc20_token_abi_str))
	bc_config := bc.Config()
	value := big.NewInt(0)
	gas_limit := uint64(50000000)
	gas_price := big.NewInt(1)
	fake_src_addr := common.HexToAddress("8999999999999999999999999999999999999998")
	fake_balance := big.NewInt(0)
	fake_balance.SetString("9999999999999999999999999999", 10)

	// Add a fake account with huge balance so we have money to pay for gas to execute instructions on the EVM
	statedb.AddBalance(fake_src_addr, fake_balance)

	// Encode input for retrieving token balance
	input, err = erc20_token_abi.Pack("decimals")
	if err != nil {
		fatalf("Can't pack input for decimals: %v", err)
	}

	// Getting token decimals
	msg := types.NewMessage(fake_src_addr, contract_address, 0, value, gas_limit, gas_price, input, false)
	evm_context := token.NewEVMContext(msg, block.Header(), bc, nil)
	evm := vm.NewEVM(evm_context, statedb, bc_config, vm_config)
	gp := new(token.GasPool).AddGas(math.MaxUint64)

	ret, gas, failed, err := token.ApplyMessage(evm, msg, gp)
	if failed {
		log.Debug(fmt.Sprintf("get_ERC20_token_decimals: vm err for symbol: %v, failed=%v", err, failed))
		return nil, fmt.Errorf("vm err")
	}

	if err != nil {
		log.Debug(fmt.Sprintf("get_ERC20_token_decimals: getting 'balanceOf' caused error in vm: %v", err))
		return nil, err
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("len(ret)==0")
	}

	decimals := new(uint8)
	if !((err != nil) || (failed)) {
		err = erc20_token_abi.Unpack(&decimals, "decimals", ret)
		if err != nil {
			fatalf("Can't upack decimals output from the EVM: %v", err)
		}
	}

	_ = gas

	return decimals, nil
}

func WriteTokenBalances(extdb ExtDB, bc BlockChain, block *types.Block, receipts types.Receipts) error {
	blockchain_config := bc.Config()
	statedb, err := bc.StateAt(block.Root())
	if err != nil {
		log.Error("Can't get StateAt()", "block_num", block.NumberU64())

		return err
	}

	if blockchain_config.DAOForkSupport && blockchain_config.DAOForkBlock != nil && blockchain_config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}

	for _, receipt := range receipts {
		if receipt.ContractAddress != (common.Address{}) { // contract creation
			tokenBalanceContract := new(exttypes.TokenBalance)
			tokenBalanceContract.TokenAddress = receipt.ContractAddress
			tokenBalanceContract.HolderAddress = receipt.ContractAddress
			tokenBalanceContract.BlockNumber = block.Number()
			tokenBalanceContract.BlockHash = block.Hash()
			tokenBalanceContract.TokenDecimals, err = get_ERC20_token_decimals_from_EVM(bc, statedb, block, &tokenBalanceContract.TokenAddress)
			tokenBalanceContract.HolderBalance, err = get_ERC20_token_balance_from_EVM(bc, statedb, block, &tokenBalanceContract.TokenAddress, &tokenBalanceContract.HolderAddress)
			if err == nil {
				extdb.WriteTokenBalance(tokenBalanceContract)
			}
			for _, transaction := range block.Transactions() {
				if transaction.Hash() == receipt.TxHash {
					signer := types.MakeSigner(bc.Config(), block.Number())
					from, _ := types.Sender(signer, transaction)

					tokenBalanceOwner := new(exttypes.TokenBalance)
					tokenBalanceOwner.TokenAddress = receipt.ContractAddress
					tokenBalanceOwner.HolderAddress = from
					tokenBalanceOwner.BlockNumber = block.Number()
					tokenBalanceOwner.BlockHash = block.Hash()
					tokenBalanceOwner.TokenDecimals = tokenBalanceContract.TokenDecimals
					tokenBalanceOwner.HolderBalance, err = get_ERC20_token_balance_from_EVM(bc, statedb, block, &tokenBalanceOwner.TokenAddress, &tokenBalanceOwner.HolderAddress)
					if err == nil {
						extdb.WriteTokenBalance(tokenBalanceOwner)
					}
				}
			}
		}
		for _, event := range receipt.Logs {
			if len(event.Topics) == 0 || len(event.Topics) < 3 {
				continue
			}
			if 0 == bytes.Compare(event.Topics[0].Bytes(), erc20_transfer_event_signature) {
				tokenBalanceFrom := new(exttypes.TokenBalance)
				tokenBalanceFrom.TokenAddress = event.Address
				tokenBalanceFrom.HolderAddress = common.BytesToAddress(event.Topics[1].Bytes())
				tokenBalanceFrom.BlockNumber = block.Number()
				tokenBalanceFrom.BlockHash = event.BlockHash
				tokenBalanceFrom.TokenDecimals, err = get_ERC20_token_decimals_from_EVM(bc, statedb, block, &tokenBalanceFrom.TokenAddress)
				tokenBalanceFrom.HolderBalance, err = get_ERC20_token_balance_from_EVM(bc, statedb, block, &tokenBalanceFrom.TokenAddress, &tokenBalanceFrom.HolderAddress)
				if err == nil {
					extdb.WriteTokenBalance(tokenBalanceFrom)
				}

				tokenBalanceTo := new(exttypes.TokenBalance)
				tokenBalanceTo.TokenAddress = event.Address
				tokenBalanceTo.HolderAddress = common.BytesToAddress(event.Topics[2].Bytes())
				tokenBalanceTo.BlockNumber = block.Number()
				tokenBalanceTo.BlockHash = event.BlockHash
				tokenBalanceTo.TokenDecimals = tokenBalanceFrom.TokenDecimals
				tokenBalanceTo.HolderBalance, err = get_ERC20_token_balance_from_EVM(bc, statedb, block, &tokenBalanceTo.TokenAddress, &tokenBalanceTo.HolderAddress)
				if err == nil {
					extdb.WriteTokenBalance(tokenBalanceTo)
				}
			}
			if 0 == bytes.Compare(event.Topics[0].Bytes(), mint_event_signature) ||
				0 == bytes.Compare(event.Topics[0].Bytes(), burn_event_signature) ||
				0 == bytes.Compare(event.Topics[0].Bytes(), jntmint_event_signature) ||
				0 == bytes.Compare(event.Topics[0].Bytes(), jntburn_event_signature) {
				tokenBalanceTo := new(exttypes.TokenBalance)
				tokenBalanceTo.TokenAddress = event.Address
				tokenBalanceTo.HolderAddress = common.BytesToAddress(event.Topics[1].Bytes())
				tokenBalanceTo.BlockNumber = block.Number()
				tokenBalanceTo.BlockHash = event.BlockHash
				tokenBalanceTo.TokenDecimals, err = get_ERC20_token_decimals_from_EVM(bc, statedb, block, &tokenBalanceTo.TokenAddress)
				tokenBalanceTo.HolderBalance, err = get_ERC20_token_balance_from_EVM(bc, statedb, block, &tokenBalanceTo.TokenAddress, &tokenBalanceTo.HolderAddress)
				if err == nil {
					extdb.WriteTokenBalance(tokenBalanceTo)
				}
			}
		}
	}

	return nil
}
