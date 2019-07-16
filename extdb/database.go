package extdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type ExtDBpg struct {
	conn          *sql.DB
	writeDuration mclock.AbsTime
	isSkipConn    bool
	nodeId        string
}

//var (
//	writer_headers *kafka.Writer               = nil
//	writer_bodies *kafka.Writer                = nil
//	writer_pending_transactions *kafka.Writer  = nil
//	writer_receipts *kafka.Writer              = nil
//	writer_accounts *kafka.Writer              = nil
//	writer_rewards *kafka.Writer               = nil
//	writer_internal_transactions *kafka.Writer = nil
//	writer_hain_splits *kafka.Writer           = nil
//	writer_reorgs *kafka.Writer                = nil
//)

// example: writer_headers, _ = Configure([]string{"127.0.0.1:39092"}, "1", "headers")
func Configure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: nil,
	}
	w = kafka.NewWriter(config)
	return w, nil
}

// example: Push(context.Background(), writer_headers, nil, []byte(fieldsString))
func Push(parent context.Context, writer *kafka.Writer, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer.WriteMessages(parent, message)
}

func NewExtDBpg(dbURI string) error {
	if dbURI == "null" {
		db = nil
		log.Info("Extern DB is null, all extern db operatons will be skipped")
		return nil
	}
	dbpg := &ExtDBpg{
		conn:       nil,
		isSkipConn: false,
	}
	db = dbpg
	return dbpg.Connect(dbURI)
}

func (self *ExtDBpg) Connect(dbURI string) error {
	conn, err := sql.Open("postgres", dbURI)
	self.conn = conn
	if dbURI == "" {
		self.isSkipConn = true
	}
	if err != nil {
		log.Crit("ExtDB Error when connect to extern DB", "Error", err)
	} else {
		re := regexp.MustCompile("(//.*:)(.*)(@)")
		log.Info("ExtDB Connected to extern DB", "URI", re.ReplaceAllString(dbURI, "$1****$3"))
	}
	return err
}

func (self *ExtDBpg) IsSkipConn() bool {
	return self.isSkipConn
}

func (self *ExtDBpg) Close() error {
	return self.conn.Close()
}

func (self *ExtDBpg) WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	start := mclock.Now()
	log.Debug("ExtDB write block header", "hash", blockHash, "number", blockNumber)

	fieldsString, err := self.SerializeHeaderFields(header)
	log.Debug("ExtDB header serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	var query = "INSERT INTO headers (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB header insertion", "time", common.PrettyDuration(query_duration))

	if err != nil {
		log.Warn("ExtDB Error writing header to extern DB", "Error", err)
	}
	return err
}

func (self *ExtDBpg) WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
	log.Debug("ExtDB write block body", "hash", blockHash, "number", blockNumber)
	start := mclock.Now()
	fieldsString, err := self.SerializeBodyFields(body)

	log.Debug("ExtDB body serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	var query = "INSERT INTO bodies (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB body insertion", "time", common.PrettyDuration(query_duration))

	if err != nil {
		log.Warn("ExtDB Error writing body to extern DB", "Error", err)
	}
	return nil
}

func (self *ExtDBpg) WritePendingTransaction(txHash common.Hash, transaction *types.Transaction, is_removed bool) error {
	start := mclock.Now()
	log.Debug("ExtDB write pending transaction", "tx_hash", txHash)

	var query = `INSERT INTO pending_transactions (tx_hash, fields, removed, node_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;`
	var fieldsString = "{}"
	var err error = nil
	if transaction != nil {
		fieldsString, err = self.SerializeTransactionFields(transaction)
	}
	log.Debug("ExtDB pending transaction serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	_, err = self.conn.Exec(query, txHash.Hex(), fieldsString, is_removed, self.nodeId)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB pending transaction insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts *exttypes.ReceiptsContainer) error {
	start := mclock.Now()
	log.Debug("ExtDB write receipts", "hash", blockHash, "number", blockNumber)

	var query = "INSERT INTO receipts (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeReceiptsFields(receipts)
	log.Debug("ExtDB receipts serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	if err != nil {
		return err
	}
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB receipts insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error {
	start := mclock.Now()
	log.Debug("ExtDB write state object", "hash", blockHash, "number", blockNumber)

	var query = "INSERT INTO accounts (block_hash, block_number, address, fields) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeStateObjectFields(obj)
	log.Debug("ExtDB account serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, addr.Hex(), fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB account insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteRewards(blockHash common.Hash, blockNumber uint64, addr common.Address, blockReward *exttypes.BlockReward) error {
	start := mclock.Now()
	log.Debug("ExtDB write rewards", "hash", blockHash, "number", blockNumber, "miner", addr)

	var query = "INSERT INTO rewards (block_hash, block_number, address, fields) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	fieldsString, err := self.SerializeBlockRewardsFields(blockReward)
	log.Debug("ExtDB rewards serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, addr.Hex(), fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB rewards insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteInternalTransaction(intTransaction *exttypes.InternalTransaction) error {
	start := mclock.Now()
	log.Debug("ExtDB write internal transaction",
		"block_number", intTransaction.BlockNumber.Uint64(),
		"op", intTransaction.Operation)

	var query = `INSERT INTO internal_transactions (block_number, block_hash, parent_tx_hash, index, type, timestamp, fields)
                 VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING;`

	fieldsString, err := self.SerializeInternalTransactionFields(intTransaction)
	log.Debug("ExtDB internal transaction serialization", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	_, err = self.conn.Exec(query, intTransaction.BlockNumber.Uint64(), intTransaction.BlockHash.Hex(),
		intTransaction.ParentTxHash.Hex(), intTransaction.Index, intTransaction.Operation,
		intTransaction.TimeStamp.Uint64(), fieldsString)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB internal transaction insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteTokenBalance(tokenBalance *exttypes.TokenBalance) error {
	start := mclock.Now()
	log.Debug("ExtDB write token balance",
		"block_number", tokenBalance.BlockNumber.Uint64(),
		"block_hash", tokenBalance.BlockHash.Hex(),
		"token_address", tokenBalance.TokenAddress.Hex(),
		"holder_address", tokenBalance.HolderAddress.Hex(),
		"balance", tokenBalance.HolderBalance.Uint64())

	var query = `INSERT INTO token_holders (block_number, block_hash, token_address, holder_address, balance)
                 VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;`

	start = mclock.Now()
	_, err := self.conn.Exec(query, tokenBalance.BlockNumber.Uint64(), tokenBalance.BlockHash.Hex(),
		tokenBalance.TokenAddress.Hex(), tokenBalance.HolderAddress.Hex(), tokenBalance.HolderBalance.Uint64())
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB write token balance", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) WriteReorg(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	time.Sleep(1000 * 10)
	start := mclock.Now()
	log.Debug("ExtDB write block reorg", "hash", blockHash, "number", blockNumber)

	headerString, err := self.SerializeHeaderFields(header)
	log.Debug("ExtDB header serialization reorg", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	var query = "INSERT INTO reorgs (block_hash, block_number, header, reinserted, split_id, node_id) VALUES ($1, $2, $3, false, $4, $5) ON CONFLICT DO NOTHING;"
	if tx != nil {
		_, err = tx.Exec(query, blockHash.Hex(), blockNumber, headerString, split_id, self.nodeId)
	} else {
		_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, headerString, split_id, self.nodeId)
	}
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB reorg insertion", "time", common.PrettyDuration(query_duration))

	if err != nil {
		log.Warn("ExtDB Error writing reorg to extern DB", "Error", err)
	}
	return err
}

func (self *ExtDBpg) ReinsertBlock(tx *sql.Tx, split_id int, blockHash common.Hash, blockNumber uint64, header *types.Header) error {
	start := mclock.Now()
	log.Debug("ExtDB reinsert block", "hash", blockHash, "number", blockNumber)

	headerString, err := self.SerializeHeaderFields(header)
	log.Debug("ExtDB header serialization reinsert block", "time", common.PrettyDuration(mclock.Now()-start))
	start = mclock.Now()
	var query = "INSERT INTO reorgs (block_hash, block_number, header, reinserted, split_id, node_id) VALUES ($1, $2, $3, true, $4, $5) ON CONFLICT DO NOTHING;"
	if tx != nil {
		_, err = tx.Exec(query, blockHash.Hex(), blockNumber, headerString, split_id, self.nodeId)
	} else {
		_, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, headerString, split_id, self.nodeId)
	}
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB reinsert block insertion", "time", common.PrettyDuration(query_duration))

	if err != nil {
		log.Warn("ExtDB Error reinserting block to extern DB", "Error", err)
	}
	return err
}

func (self *ExtDBpg) WriteChainSplit(tx *sql.Tx, common_block_number uint64, common_block_hash common.Hash, drop_length int, drop_block_hash common.Hash, add_length int, add_block_hash common.Hash) (int, error) {
	start := mclock.Now()
	log.Debug("ExtDB write chain split", "common hash", common_block_hash, "common number", common_block_number, "drop length", drop_length, "drop hash", drop_block_hash, "add length", add_length)

	var query = "INSERT INTO chain_splits (common_block_number, common_block_hash, drop_length, drop_block_hash, add_length, add_block_hash, node_id) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Warn("ExtDB Error prepare write chain split query", "Error", err)
		return 0, err
	}

	defer stmt.Close()

	var chain_split_id int
	err = stmt.QueryRow(common_block_number, common_block_hash.Hex(), drop_length, drop_block_hash.Hex(), add_length, add_block_hash.Hex(), self.nodeId).Scan(&chain_split_id)

	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB chain split insertion", "time", common.PrettyDuration(query_duration))

	if err != nil {
		log.Warn("ExtDB Error writing chain split to extern DB", "Error", err)
	}
	return chain_split_id, err
}

func (self *ExtDBpg) BeginTx() (*sql.Tx, error) {
	tx, err := self.conn.BeginTx(context.Background(), nil)
	return tx, err
}

func (self *ExtDBpg) CloseTx(tx *sql.Tx, commit bool) error {
	if commit {
		err := tx.Commit()
		return err
	} else {
		err := tx.Rollback()
		return err
	}
}

func (self *ExtDBpg) ResetDbWriteDuration() error {
	self.writeDuration = mclock.AbsTime(0)
	return nil
}

func (self *ExtDBpg) UpdateDbWriteDuration(duration mclock.AbsTime) error {
	self.writeDuration += duration
	return nil
}

func (self *ExtDBpg) GetDbWriteDuration() mclock.AbsTime {
	return self.writeDuration
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

func (self *ExtDBpg) WriteChainEvent(
	block_number uint64,
	block_hash common.Hash,
	parent_block_hash common.Hash,
	event_type string,
	drop_length int,
	drop_block_hash common.Hash,
	add_length int,
	add_block_hash common.Hash) error {

	start := mclock.Now()
	log.Debug("ExtDB write chain event", "block hash", block_hash, "block number", block_number, "event type", event_type)

	var query = "INSERT INTO chain_events (block_number, block_hash, parent_block_hash, type, common_block_number, common_block_hash, drop_length, drop_block_hash, add_length, add_block_hash, node_id) VALUES ($1, $2, $3, $4, 0, '', $5, $6, $7, $8, $9);"
	_, err := self.conn.Exec(query, block_number, block_hash.Hex(), parent_block_hash.Hex(), event_type, drop_length, drop_block_hash.Hex(), add_length, add_block_hash.Hex(), self.nodeId)
	query_duration := mclock.Now() - start
	self.UpdateDbWriteDuration(query_duration)
	log.Debug("ExtDB chain event insertion", "time", common.PrettyDuration(query_duration))

	return err
}

func (self *ExtDBpg) SetNodeId(nodeId enode.ID) error {
	self.nodeId = common.ToHex(nodeId[:])
	return nil
}
