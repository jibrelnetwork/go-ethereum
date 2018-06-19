package extdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

type ExtDB interface {
	Connect(dbURI string) error
	Close() error
	WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error
	WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error
	WritePendingTransaction(txHash common.Hash, transaction *types.Transaction) error
	WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error
	WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error
	WriteRewards(blockHash common.Hash, blockNumber uint64, addr common.Address, blockReward *exttypes.BlockReward) error
	WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error
	NewBlockNotify(blockNumber uint64) error
}

var (
	ExtDbUriFlag = cli.StringFlag{
		Name:  "extdb",
		Usage: "Extern DB connection string",
	}

	db ExtDB
)

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	if db != nil {
		return db.WriteBlockHeader(blockHash, blockNumber, header)
	}
	return nil
}

func WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
	if db != nil {
		return db.WriteBlockBody(blockHash, blockNumber, body)
	}
	return nil
}

func WritePendingTransaction(txHash common.Hash, transaction *types.Transaction) error {
	if db != nil {
		return db.WritePendingTransaction(txHash, transaction)
	}
	return nil
}

func WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error {
	if db != nil {
		return db.WriteReceipts(blockHash, blockNumber, receipts)
	}
	return nil
}

func WriteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address, dumpAccount interface{}) error {
	if db != nil {
		if err := db.WriteStateObject(blockHash, blockNumber, address, dumpAccount); err != nil {
			return err
		}
	}
	return nil
}

func DeleteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address) error {
	if db != nil {
		log.Debug("Stubbed delete state object in ext db", "Addr", address.Hex())
	}
	return nil

}

func WriteRewards(blockHash common.Hash, blockNumber uint64, address common.Address, blockReward *exttypes.BlockReward) error {
	if db != nil {
		if err := db.WriteRewards(blockHash, blockNumber, address, blockReward); err != nil {
			return err
		}
	}
	return nil
}

func WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error {
	if db != nil {
		if err := db.WriteInternalTransaction(intTransaction); err != nil {
			return err
		}
	}
	return nil
}

func NewBlockNotify(blockNumber uint64) error {
	if db != nil {
		if err := db.NewBlockNotify(blockNumber); err != nil {
			return err
		}
	}
	return nil
}
