// Package exttypes contains data types related to Token Balances
package exttypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TokenBalance struct {
	TokenAddress  common.Address
	HolderAddress common.Address
	HolderBalance *big.Int
	BlockNumber   *big.Int
	BlockHash     common.Hash
}
