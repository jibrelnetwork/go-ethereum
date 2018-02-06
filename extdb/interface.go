package extdb

import (
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/core/types"
    "gopkg.in/urfave/cli.v1"
)


type ExtDB interface {

    Connect(dbURI string)                                                                 error
    
    WriteBlockHeader(blockHash common.Hash, blockNumber uint64, header *types.Header)     error
    WriteBlockBody(blockHash common.Hash, blockNumber uint64, body *types.Body)           error
    WriteUncles(blockHash common.Hash, blockNumber uint64, uncles []*types.Header)        error
    WriteTransaction(blockHash common.Hash, blockNumber uint64, index int, transaction *types.Transaction) error
    WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions)     error
    WriteReceipts(blockHash common.Hash, blockNumber uint64, receipts types.Receipts)     error
    WriteStateObject(blockHash common.Hash, blockNumber uint64, addr common.Address, obj interface{}) error
}


var (	
	ExtDbUriFlag = cli.StringFlag{
        Name:  "extdb",
        Usage: "Extern DB connection string",
    }

	db ExtDB
)


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


func WriteTransaction(blockHash common.Hash, blockNumber uint64, index int, transaction *types.Transaction) error {
    if db != nil{
        return db.WriteTransaction(blockHash, blockNumber, index, transaction)
    }
    return nil
}


func WriteTransactions(blockHash common.Hash, blockNumber uint64, transactions types.Transactions) error {
    if db != nil{
        return db.WriteTransactions(blockHash, blockNumber, transactions)
    }
    return nil
}


func WriteUncles(blockHash common.Hash, blockNumber uint64, uncles []*types.Header) error {
    if db != nil{
        return db.WriteUncles(blockHash, blockNumber, uncles)
    }
    return nil
}


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
