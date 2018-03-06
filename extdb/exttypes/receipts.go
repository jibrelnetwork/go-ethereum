// Package exttypes contains data types related to Receipts
package exttypes

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiptsContainer struct {
	Receipts []*types.Receipt
}
