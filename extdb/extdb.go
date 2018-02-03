package extdb

import (
    // "fmt"
    "database/sql"
    "encoding/json"
    "regexp"

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
    // WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions)                      error
    WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts)                        error
    WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error
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
        re := regexp.MustCompile("(//.*:)(.*)(@)")
        log.Info("Connected to extern DB", "URI", re.ReplaceAllString(dbURI, "$1****$3"))
    }
    return err
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
    // var fieldsString = "ASD"
    // var err = nil
    var query = "INSERT INTO blocks (block_hash, block_number, fields) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
    _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, fieldsString)
    
    if err != nil {
        log.Warn("Error writing header to extern DB", "Error", err)
    }
    return nil

}


// func (self *ExtDBpg) WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions) error {
//     var query = "INSERT INTO transactions (block_hash, block_number, tx_hash, index, fields) VALUES ($1, $2, $3, $4, $5)"
//     for i, tx := range transactions {
//         fieldsString, err := self.SerializeTransactionFields(tx)
//         _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, tx.Hash().Hex(), i, fieldsString)
//         if err != nil {
//             return err
//         }
//     }
//     return nil
// }


func (self *ExtDBpg) WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts) error {
    var query = "INSERT INTO receipts (block_hash, block_number, tx_hash, index, fields) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"
    for i, receipt := range receipts {
        fieldsString, err := self.SerializeReceiptFields(receipt)
        _, err = self.conn.Exec(query, blockHash.Hex(), blockNumber, receipt.TxHash.Hex(), i, fieldsString)
        if err != nil {
            return err
        }
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

func (self *ExtDBpg) SerializeStateObjectFields(dumpAccount interface{}) (string, error) {
    b, err := json.Marshal(dumpAccount)
    return string(b), err
}


// func (self *ExtDBpg) SerializeTransactionFields(transaction *types.Transaction) (string, error) {
//     b, err := json.Marshal(transaction)
//     return string(b), err
// }


func WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header) error {
    if db != nil{
        return db.WriteBlockHeader(blockHash, blockNumber, header)
    }
    return nil
}


func WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body) error {
    if db != nil{
        return db.WriteBlockBody(blockHash, blockNumber, body)
    }
    return nil
}


// func WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions) error {
//     return db.WriteTransactions(blockHash, blockNumber, transactions)
// }


func WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts) error {
    if db != nil{
        return db.WriteReceipts(blockHash, blockNumber, receipts)
    }
    return nil
}


func WriteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address, dumpAccount interface{}) error {

    if db != nil{
        if err := db.WriteStateObject(blockHash, blockNumber, address, dumpAccount); err != nil {
            return err
        }
    }
    return nil

}


func DeleteStateObject(blockHash common.Hash, blockNumber uint64, address common.Address) error {

    if db != nil{
        log.Info("Stubbed delete state object in ext db", "Addr", address.Hex())
    }
    return nil

}
