// Package exttypes contains data types related to Internal Transaction
package exttypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type InternalTransaction struct {
	BlockNumber  *big.Int
	TimeStamp    *big.Int
	From         *common.Address
	To           *common.Address
	Value        *big.Int
	GasUsed      uint64
	Status       string
	ParentTxHash common.Hash
}
