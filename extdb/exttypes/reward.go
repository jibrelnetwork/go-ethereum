// Package exttypes contains data types related to Block Reward
package exttypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type UncleReward struct {
	Miner         common.Address
	UnclePosition int
	BlockReward   *big.Int
}

type BlockReward struct {
	BlockNumber          *big.Int
	TimeStamp            *big.Int
	BlockMiner           common.Address
	BlockReward          *big.Int
	Uncles               []*UncleReward
	UncleInclusionReward *big.Int
}
