package extdb

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"gopkg.in/urfave/cli.v1"
)

type ExtDB interface {
	Connect(dbURI string) error
	Close() error
	WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error
	WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error
	WritePendingTransaction(txHash common.Hash, transaction *types.Transaction, is_removed bool) error
	WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error
	WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error
	WriteRewards(blockHash common.Hash, blockNumber uint64, addr common.Address, blockReward *exttypes.BlockReward) error
	WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error
	WriteReorg(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error
	WriteChainSplit(tx *sql.Tx, common_block_number uint64, common_block_hash common.Hash, drop_length int, drop_block_hash common.Hash, add_length int, add_block_hash common.Hash) (int, error)
	ReinsertBlock(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error
	BeginTx() (*sql.Tx, error)
	CloseTx(tx *sql.Tx, commit bool) error
	GetDbWriteDuration() mclock.AbsTime
	ResetDbWriteDuration() error
	IsSkipConn() bool
	WriteChainEvent(
		block_number uint64,
		block_hash common.Hash,
		event_type string,
		common_block_number uint64,
		common_block_hash common.Hash,
		drop_length int,
		drop_block_hash common.Hash,
		add_length int,
		add_block_hash common.Hash) error
	SetNodeId(nodeId enode.ID) error
}

var (
	ExtDbUriFlag = cli.StringFlag{
		Name:  "extdb",
		Usage: "Extern DB connection string",
	}

	db ExtDB
)

func Close() error {
	if db != nil && !db.IsSkipConn() {
		return db.Close()
	}
	return nil
}

func WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	if db != nil && !db.IsSkipConn() {
		return db.WriteBlockHeader(blockHash, blockNumber, header)
	}
	return nil
}

func WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
	if db != nil && !db.IsSkipConn() {
		return db.WriteBlockBody(blockHash, blockNumber, body)
	}
	return nil
}

func WritePendingTransaction(txHash common.Hash, transaction *types.Transaction, is_removed bool) error {
	if db != nil && !db.IsSkipConn() {
		return db.WritePendingTransaction(txHash, transaction, is_removed)
	}
	return nil
}

func WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error {
	if db != nil && !db.IsSkipConn() {
		return db.WriteReceipts(blockHash, blockNumber, receipts)
	}
	return nil
}

func WriteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address, dumpAccount interface{}) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.WriteStateObject(blockHash, blockNumber, address, dumpAccount); err != nil {
			return err
		}
	}
	return nil
}

func DeleteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address) error {
	if db != nil && !db.IsSkipConn() {
		log.Debug("Stubbed delete state object in ext db", "Addr", address.Hex())
	}
	return nil

}

func WriteRewards(blockHash common.Hash, blockNumber uint64, address common.Address, blockReward *exttypes.BlockReward) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.WriteRewards(blockHash, blockNumber, address, blockReward); err != nil {
			return err
		}
	}
	return nil
}

func WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.WriteInternalTransaction(intTransaction); err != nil {
			return err
		}
	}
	return nil
}

func WriteReorg(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.WriteReorg(tx, split_id, blockHash, blockNumber, header); err != nil {
			return err
		}
	}
	return nil
}

func WriteChainSplit(tx *sql.Tx, common_block_number uint64, common_block_hash common.Hash, drop_length int, drop_block_hash common.Hash, add_length int, add_block_hash common.Hash) (int, error) {
	if db != nil && !db.IsSkipConn() {
		return db.WriteChainSplit(tx, common_block_number, common_block_hash, drop_length, drop_block_hash, add_length, add_block_hash)
	}
	return 0, nil
}

func BeginTx() (*sql.Tx, error) {
	if db != nil && !db.IsSkipConn() {
		return db.BeginTx()
	}
	return nil, nil
}

func CloseTx(tx *sql.Tx, commit bool) error {
	if db != nil && !db.IsSkipConn() {
		return db.CloseTx(tx, commit)
	}
	return nil
}

func ReinsertBlock(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.ReinsertBlock(tx, split_id, blockHash, blockNumber, header); err != nil {
			return err
		}
	}
	return nil
}

func ResetDbWriteDuration() error {
	if db != nil && !db.IsSkipConn() {
		db.ResetDbWriteDuration()
	}
	return nil
}

func GetDbWriteDuration() mclock.AbsTime {
	if db != nil && !db.IsSkipConn() {
		return db.GetDbWriteDuration()
	}
	return mclock.AbsTime(0)
}

func WriteChainEvent(
	block_number uint64,
	block_hash common.Hash,
	event_type string,
	common_block_number uint64,
	common_block_hash common.Hash,
	drop_length int,
	drop_block_hash common.Hash,
	add_length int,
	add_block_hash common.Hash) error {
	if db != nil && !db.IsSkipConn() {
		if err := db.WriteChainEvent(block_number, block_hash, event_type, common_block_number, common_block_hash, drop_length, drop_block_hash, add_length, add_block_hash); err != nil {
			return err
		}
	}
	return nil
}

func SetNodeId(nodeId enode.ID) error {
	if db != nil {
		return db.SetNodeId(nodeId)
	}
	return nil
}
