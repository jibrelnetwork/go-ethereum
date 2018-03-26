// Package exttypes contains data types related to Internal Transaction
package exttypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type InternalTransaction struct {
	BlockNumber  *big.Int
	Operation    string
	CallDepth    int
	TimeStamp    *big.Int
	From         *common.Address
	To           *common.Address
	Value        *big.Int
	GasLimit     uint64
	Status       string
	ParentTxHash common.Hash
}
