package extdb

import (
	"database/sql"
	"encoding/json"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/log"
	_ "github.com/lib/pq"
)

type ExtDBpg struct {
	conn *sql.DB
}

func NewExtDBpg(dbURI string) error {
	dbpg := &ExtDBpg{
		conn: nil,
	}
	db = dbpg
	return dbpg.Connect(dbURI)
}

func (self *ExtDBpg) Connect(dbURI string) error {
	conn, err := sql.Open("postgres", dbURI)
	self.conn = conn
	if err != nil {
		log.Crit("Error when connect to extern DB", "Error", err)
	} else {
		re := regexp.MustCompile("(//.*:)(.*)(@)")
		log.Info("Connected to extern DB", "URI", re.ReplaceAllString(dbURI, "$1****$3"))
	}
	return err
}

func (self *ExtDBpg) Close() error {
	return self.conn.Close()
}

func (self *ExtDBpg) WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	log.Info("ExtDB write block header", "hash", blockHash, "number", blockNumber)

	fieldsString, err := self.SerializeHeaderFields(header)
	var query = "INSERT INTO headers (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)

	if err != nil {
		log.Warn("Error writing header to extern DB", "Error", err)
	}
	return err
}

func (self *ExtDBpg) WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
	log.Info("ExtDB write block body", "hash", blockHash, "number", blockNumber)

	fieldsString, err := self.SerializeBodyFields(body)
	var query = "INSERT INTO bodies (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)

	if err != nil {
		log.Warn("Error writing body to extern DB", "Error", err)
	}
	return nil
}

func (self *ExtDBpg) WritePendingTransaction(txHash common.Hash, transaction *types.Transaction) error {
	log.Info("ExtDB write pendinf transaction", "tx_hash", txHash)

	var query = `INSERT INTO pending_transactions (tx_hash, fields)
                 VALUES ($1, $2)
                 ON CONFLICT (tx_hash) DO UPDATE
                 SET fields=excluded.fields;`

	fieldsString, err := self.SerializeTransactionFields(transaction)
	_, err = self.conn.Exec(query, txHash.Hex(), fieldsString)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error {
	log.Info("ExtDB write receipts", "hash", blockHash, "number", blockNumber)

	var query = "INSERT INTO receipts (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeReceiptsFields(receipts)
	if err != nil {
		return err
	}
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error {
	var query = "INSERT INTO accounts (block_hash, block_number, address, fields) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeStateObjectFields(obj)
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, addr.Hex(), fieldsString)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) WriteRewards(blockHash common.Hash, blockNumber uint64, addr common.Address, blockReward *exttypes.BlockReward) error {
	log.Info("ExtDB write reewards", "hash", blockHash, "number", blockNumber, "miner", addr)

	var query = "INSERT INTO rewards (block_hash, block_number, address, fields) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeBlockRewardsFields(blockReward)
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, addr.Hex(), fieldsString)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error {
	log.Info("ExtDB write internal transaction",
		"block_number", intTransaction.BlockNumber.Uint64(),
		"op", intTransaction.Operation)

	var query = `INSERT INTO internal_transactions (block_number, type, timestamp, fields)
                 VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;`

	fieldsString, err := self.SerializeInternalTransactionFields(intTransaction)
	_, err = self.conn.Exec(query, intTransaction.BlockNumber.Uint64(), intTransaction.Operation, intTransaction.TimeStamp.Uint64(), fieldsString)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) NewBlockNotify(blockNumber uint64) error {
	var query = `select pg_notify('newblock', CAST($1 AS text));`

	_, err := self.conn.Exec(query, blockNumber)
	if err != nil {
		return err
	}
	return nil
}

func (self *ExtDBpg) SerializeHeaderFields(header *types.Header) (string, error) {
	b, err := json.Marshal(header)
	return string(b), err
}

func (self *ExtDBpg) SerializeBodyFields(body *types.Body) (string, error) {
	b, err := json.Marshal(body)
	return string(b), err
}

func (self *ExtDBpg) SerializeReceiptsFields(receipts *exttypes.ReceiptsContainer) (string, error) {
	b, err := json.Marshal(receipts)
	return string(b), err
}

func (self *ExtDBpg) SerializeStateObjectFields(dumpAccount interface{}) (string, error) {
	b, err := json.Marshal(dumpAccount)
	return string(b), err
}

func (self *ExtDBpg) SerializeTransactionFields(transaction *types.Transaction) (string, error) {
	b, err := json.Marshal(transaction)
	return string(b), err
}

func (self *ExtDBpg) SerializeUncleFields(uncle *types.Header) (string, error) {
	b, err := json.Marshal(uncle)
	return string(b), err
}

func (self *ExtDBpg) SerializeBlockRewardsFields(blockReward *exttypes.BlockReward) (string, error) {
	b, err := json.Marshal(blockReward)
	return string(b), err
}

func (self *ExtDBpg) SerializeInternalTransactionFields(intTransaction *exttypes.InternalTransaction) (string, error) {
	b, err := json.Marshal(intTransaction)
	return string(b), err
}
