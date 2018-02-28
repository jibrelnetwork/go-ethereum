package tests

import (
	"fmt"
	_"fmt"
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/extdb"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	_ "github.com/lib/pq"
)

var (
	extdb_constr = flag.String("extdb", "", "Extern DB connection string")
)


func readBlocks(t *testing.T, db *sql.DB) {
	var (
		testdb, _ = ethdb.NewMemDatabase()
		gspec     = &core.Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		_, _ = core.GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), testdb, 1, nil)

		prevHeader, curHeader *types.Header
		curBody *types.Body
		curReceipts *exttypes.ReceiptsContainer
		blockNumber int64
		blockHash string
		headersFields, bodiesFields, receiptsFields []byte
	)

	engine := ethash.NewFaker()
	chain, _ := core.NewBlockChain(testdb, nil, params.TestChainConfig, engine, vm.Config{})
	
	rows, err := db.Query(`SELECT h.block_number, h.block_hash, h.fields, b.fields, r.fields
		FROM headers AS h 
		LEFT JOIN bodies AS b ON b.block_number=h.block_number
		LEFT JOIN receipts AS r ON r.block_number=h.block_number
		ORDER BY h.block_number;`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&blockNumber, &blockHash, &headersFields, &bodiesFields, &receiptsFields)
		if err != nil {
			t.Fatal(err)
		}
		curHeader = new(types.Header)
		curBody = new(types.Body)
		curReceipts = new(exttypes.ReceiptsContainer)

		err = json.Unmarshal(headersFields, curHeader)
		if err !=nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(bodiesFields, curBody)
		if err !=nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(receiptsFields, curReceipts)
		if err !=nil {
			fmt.Printf("%s", receiptsFields)
			t.Fatal(err)
		}

		handleBlock(t, chain, engine, blockNumber, blockHash, curHeader, prevHeader, curBody)
		handleReceipts(t, blockNumber, curReceipts, curHeader)

		prevHeader = curHeader
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}


func readAccounts(t *testing.T, db *sql.DB, addrs []common.Address, statedb *state.StateDB) {
	var (
		blockNumber int64
		blockHash, address string
		fields []byte
		stateDump *state.DumpAccount
	)

	rows, err := db.Query(`SELECT a2.block_number, a2.block_hash, a2.address, a2.fields FROM 
		(SELECT DISTINCT address, MAX(block_number) AS block_number FROM accounts GROUP BY address) AS a1
		LEFT JOIN accounts AS a2 ON a2.block_number=a1.block_number AND a2.address=a1.address;`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&blockNumber, &blockHash, &address, &fields)
		if err != nil {
			t.Fatal(err)
		}
		stateDump = new(state.DumpAccount)
		err = json.Unmarshal(fields, stateDump)
		if err !=nil {
			t.Fatal(err)
		}
		handleAccount(t, blockNumber, blockHash, address, stateDump, addrs, statedb)
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}


func handleBlock(t *testing.T, chain *core.BlockChain, engine *ethash.Ethash, blockNumber int64, blockHash string, curHeader *types.Header, prevHeader *types.Header, curBody *types.Body) {
	if blockNumber != curHeader.Number.Int64() {
		t.Fatalf("block_number mismatch: have %x, want %x", blockNumber, curHeader.Number)
	}
	if blockHash != curHeader.Hash().Hex() {
		t.Fatalf("block_hash mismatch: have %x, want %x", blockHash, curHeader.Hash())
	}
	if prevHeader != nil {
		if prevHeader.Hash() != curHeader.ParentHash {
			t.Fatalf("parent_hash mismatch: have %x, want %x", prevHeader.Hash(), curHeader.ParentHash)
		}
		err := engine.VerifyHeader2(chain, curHeader, prevHeader, false, true)
		if err != nil {
			t.Fatal(err)
		}
	}
	if hash := types.CalcUncleHash(curBody.Uncles); hash != curHeader.UncleHash {
		t.Fatalf("uncle root hash mismatch: have %x, want %x", hash, curHeader.UncleHash)
	}
	if hash := types.DeriveSha(types.Transactions(curBody.Transactions)); hash != curHeader.TxHash {
		t.Fatalf("transaction root hash mismatch: have %x, want %x", hash, curHeader.TxHash)
	}
}


func handleAccount(t *testing.T, blockNumber int64, blockHash string, address string, stateDump *state.DumpAccount, addrs []common.Address, statedb *state.StateDB) {
	var addr common.Address = common.BytesToAddress(common.FromHex(address))

	stateBalance := statedb.GetBalance(addr)
	if stateDump.Balance != stateBalance.String() {
		t.Fatalf("state_balance mismatch: have %x, want %x, address %s", stateBalance, stateDump.Balance, address)
	}
}


func handleReceipts(t *testing.T, blockNumber int64, receipts *exttypes.ReceiptsContainer, header *types.Header) {
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts.Receipts)
	if rbloom != header.Bloom {
		t.Fatalf("invalid bloom (remote: %x  local: %x), block_number: %x", header.Bloom, rbloom, blockNumber)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveSha(types.Receipts(receipts.Receipts))
	if receiptSha != header.ReceiptHash {
		t.Fatalf("invalid receipt root hash (remote: %x local: %x), block_number: %x", header.ReceiptHash, receiptSha, blockNumber)
	}
}


func createTestTables(t *testing.T, connectionString string) {
	db, dbErr := sql.Open("postgres", connectionString)
	if dbErr != nil {
		t.Fatalf("Filed to open database. %s", dbErr.Error())
	}
	defer db.Close()

	file, err := os.Open("../extdb/migrations/001-initial.sql")
	if err != nil {
		t.Fatalf("Open sql file failed. %s", err.Error())
	}
	
	reader := bufio.NewReader(file)
	var query string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = line[:len(line)-1]
		if line == "" {
			_, err := db.Exec(query)
			if err != nil {
				t.Fatalf("Exec sql query failed. %s", err.Error())
			}
			query = ""
		}
		query += line
	}

	file.Close()
}


func createTestDatabase(t *testing.T, connectionString string) (string, func()) {
	u, err := url.Parse(connectionString)
    if err != nil {
        t.Fatalf("Failed to parse connection string. %s", err.Error())
	}

	db, dbErr := sql.Open("postgres", connectionString)
	if dbErr != nil {
		t.Fatalf("Filed to open database. %s", dbErr.Error())
	}
  
	rand.Seed(time.Now().UnixNano())
	dbName := "jsearch" + strconv.FormatInt(rand.Int63(), 10)

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		t.Fatalf("Failed to create database. %s", err.Error())
	}

	connectionString = u.Scheme + "://" + u.User.String() + "@" + u.Host + "/" + dbName + "?" + u.RawQuery
	createTestTables(t, connectionString)

	return connectionString, func() {
		_, err := db.Exec("DROP DATABASE " + dbName)
		if err != nil {
			t.Fatalf("Drop database failed. %s", err.Error())
		}
		db.Close()
	}
}


func generateBlockchain(t *testing.T) ([]common.Address, *state.StateDB) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addrs = []common.Address{
			crypto.PubkeyToAddress(key1.PublicKey),
			crypto.PubkeyToAddress(key2.PublicKey),
			crypto.PubkeyToAddress(key3.PublicKey),
		}
		testdb, _   = ethdb.NewMemDatabase()
	)

	// Ensure that key1 has some funds in the genesis block.
	gspec := &core.Genesis{
		Config: params.TestChainConfig, //ChainConfig{HomesteadBlock: new(big.Int)},
		Alloc:  core.GenesisAlloc{addrs[0]: {Balance: big.NewInt(1000000)}},
	}
	genesis := gspec.MustCommit(testdb)

	// This call generates a chain of 5 blocks. The function runs for
	// each block and adds different features to gen based on the
	// block index.
	signer := types.HomesteadSigner{}
	blocks, _ := core.GenerateChain(gspec.Config, genesis, ethash.NewFaker(), testdb, 5, func(i int, gen *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, addr1 sends addr2 some ether.
			tx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[0]), addrs[1], big.NewInt(10000), params.TxGas, nil, nil), signer, key1)
			gen.SetCoinbase(addrs[2])
			gen.SetExtra([]byte("addr3"))
			gen.AddTx(tx)
		case 1:
			// In block 2, addr1 sends some more ether to addr2.
			// addr2 passes it on to addr3.
			tx1, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[0]), addrs[1], big.NewInt(1000), params.TxGas, nil, nil), signer, key1)
			tx2, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[1]), addrs[2], big.NewInt(1000), params.TxGas, nil, nil), signer, key2)
			gen.SetCoinbase(addrs[0])
			gen.SetExtra([]byte("addr1"))
			gen.AddTx(tx1)
			gen.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by addr3.
			gen.SetCoinbase(addrs[2])
			gen.SetExtra([]byte("addr3"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := gen.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			gen.AddUncle(b2)
			b3 := gen.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			gen.AddUncle(b3)
		}
	})

	// Import the chain. This runs all block validation rules.
	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(testdb, nil, gspec.Config, engine, vm.Config{})

	if i, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("insert error (block %d): %v\n", blocks[i].NumberU64(), err)
	}

	state, _ := blockchain.State()

	blockchain.Stop()

	return addrs, state
}


// Tests that blockchain saving works.
func TestBlockchainSaving(t *testing.T) {
	var (
		err error
		db *sql.DB
	)

	connectionString, dropDb := createTestDatabase(t, *extdb_constr)
	err = extdb.NewExtDBpg(connectionString)
	if err != nil {
		t.Fatalf("Filed to open database. %s", err.Error())
	}

	addresses, state := generateBlockchain(t)

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Filed to open database. %s", err.Error())
	}
	readBlocks(t, db)
	readAccounts(t, db, addresses, state)
	db.Close()

	extdb.Close()
	dropDb()
}
