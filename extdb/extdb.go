package extdb

import (
    // "fmt"
    "database/sql"
    "encoding/json"

    _ "github.com/lib/pq"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/core/types"
    "gopkg.in/urfave/cli.v1"
)


//connStr := "postgres://dbuser:@localhost/jsearch"

type ExtDB interface {

    Connect(dbURI string)                                                                 error
    
    WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header)     error
    WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body)           error
    WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions)                      error
    WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts)                        error
    // WriteTxLogs(blockHash common.Hash, logs []*types.Log)                                 error
    // WriteUncles(blockHash common.Hash, uncles []*types.Header)                            error


}

var ExtDbUriFlag = cli.StringFlag{
        Name:  "extdb",
        Usage: "Extern DB connection string",
    }


var db ExtDB 


type ExtDBpg struct {
    conn    * sql.DB
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
    }else{
        log.Info("Connected to extern DB", "URI", dbURI)
    }
    return err
}



func (self *ExtDBpg) WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
    log.Info("ExtDB write block header", "hash", blockHash, "number", blockNumber)

    fieldsString, err := self.SerializeHeaderFields(header)
    var query = "INSERT INTO block_headers (block_hash, block_number, fields) VALUES ($1, $2, $3)"
    _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
    
    if err != nil {
        log.Warn("Error writing header to extern DB", "Error", err)
    }
    return err
}


func (self *ExtDBpg) WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
    log.Info("ExtDB write block body", "hash", blockHash, "number", blockNumber)

    fieldsString, err := self.SerializeBodyFields(body)
    var query = "INSERT INTO block_bodies (block_hash, block_number, fields) VALUES ($1, $2, $3)"
    _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
    
    if err != nil {
        log.Warn("Error writing header to extern DB", "Error", err)
    }
    return nil

}


func (self *ExtDBpg) WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions) error {
    var query = "INSERT INTO transactions (block_hash, block_number, tx_hash, index, fields) VALUES ($1, $2, $3, $4, $5)"
    for i, tx := range transactions {
        fieldsString, err := self.SerializeTransactionFields(tx)
        _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, tx.Hash().Hex(), i, fieldsString)
        if err != nil {
            return err
        }
    }
    return nil
}


func (self *ExtDBpg) WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts) error {
    var query = "INSERT INTO receipts (block_hash, block_number, tx_hash, index, fields) VALUES ($1, $2, $3, $4, $5)"
    for i, receipt := range receipts {
        fieldsString, err := self.SerializeReceiptFields(receipt)
        _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, receipt.TxHash.Hex(), i, fieldsString)
        if err != nil {
            return err
        }
    }
    return nil
}


func (self *ExtDBpg) SerializeHeaderFields(header *types.Header) (string, error) {
    b, err := json.Marshal(header)
    // b, err := header.MarshalJSON()
    return string(b), err
}


func (self *ExtDBpg) SerializeBodyFields(body *types.Body) (string, error) {
    b, err := json.Marshal(body)
    return string(b), err
}


func (self *ExtDBpg) SerializeReceiptFields(receipt *types.Receipt) (string, error) {
    b, err := json.Marshal(receipt)
    return string(b), err
}


func (self *ExtDBpg) SerializeTransactionFields(transaction *types.Transaction) (string, error) {
    b, err := json.Marshal(transaction)
    return string(b), err
}


func WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
    return db.WriteBlockHeader(blockHash, blockNumber, header)
}


func WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
    return db.WriteBlockBody(blockHash, blockNumber, body)
}


func WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions) error {
    return db.WriteTransactions(blockHash, blockNumber, transactions)
}


func WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts) error {
    return db.WriteReceipts(blockHash, blockNumber, receipts)
}