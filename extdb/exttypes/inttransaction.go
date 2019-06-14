// Package exttypes contains data types related to Internal Transaction
package exttypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type InternalTransaction struct {
	BlockNumber  *big.Int
	BlockHash    common.Hash
	Operation    string
	CallDepth    int
	TimeStamp    *big.Int
	TxOrigin     *common.Address
	From         *common.Address
	To           *common.Address
	Value        *big.Int
	GasLimit     uint64
	Status       string
	ParentTxHash common.Hash
	Index        int
	Payload      hexutil.Bytes
}
