package tests

import (
	_ "bufio"
	"database/sql"
	"encoding/json"
	"flag"
	_ "io"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	_ "os/exec"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/extdb"
	"github.com/ethereum/go-ethereum/extdb/exttypes"
	"github.com/ethereum/go-ethereum/params"
)

var (
	extdb_constr     = flag.String("extdb", "postgresql://postgres@localhost:5432?sslmode=disable", "Extern DB connection string")
	extdb_nocreatedb = flag.Bool("nocreatedb", false, "Specifying nocreatedb will deny to create database")
	extdb_nodropdb   = flag.Bool("nodropdb", false, "Specifying nodropdb will deny to drop database")
)

func readBlocks(t *testing.T, db *sql.DB, blockchain *core.BlockChain) {
	var (
		testdb  = ethdb.NewMemDatabase()
		gspec   = &core.Genesis{Config: params.TestChainConfig}
		genesis = gspec.MustCommit(testdb)
		_, _    = core.GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), testdb, 1, nil)

		prevHeader                                  *types.Header
		blockNumber                                 int64
		blockHash                                   string
		headersFields, bodiesFields, receiptsFields []byte
	)

	engine := ethash.NewFaker()
	chain, _ := core.NewBlockChain(testdb, nil, params.TestChainConfig, engine, vm.Config{}, nil)

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

		curHeader := new(types.Header)
		err = json.Unmarshal(headersFields, curHeader)
		if err != nil {
			t.Fatal(err)
		}

		curBody := new(types.Body)
		err = json.Unmarshal(bodiesFields, curBody)
		if err != nil {
			t.Fatal(err)
		}

		curReceipts := new(exttypes.ReceiptsContainer)
		err = json.Unmarshal(receiptsFields, curReceipts)
		if err != nil {
			t.Fatal(err)
		}

		origBlock := blockchain.GetBlockByNumber(uint64(blockNumber))
		handleBlock(t, chain, engine, blockNumber, blockHash, curHeader, prevHeader, curBody, origBlock)
		handleReceipts(t, blockNumber, types.Receipts(curReceipts.Receipts), curHeader)

		prevHeader = curHeader
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}

func readAccounts(t *testing.T, db *sql.DB, origBlockchain *core.BlockChain) {
	var (
		blockNumber        int64
		blockHash, address string
		fields             []byte
		lastStateDump      *state.DumpAccount
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
		lastStateDump = new(state.DumpAccount)
		err = json.Unmarshal(fields, lastStateDump)
		if err != nil {
			t.Fatal(err)
		}
		origLastState, _ := origBlockchain.State()
		handleAccount(t, blockNumber, blockHash, address, lastStateDump, origLastState)
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}

// fmtJSON returns a pretty-printed JSON form for x.
func fmtJSON(x interface{}) string {
	js, _ := json.MarshalIndent(x, "", "\t")
	return string(js)
}

func equal(a, b interface{}) bool {
	return fmtJSON(a) != fmtJSON(b) // ignore unexported fields
}

func handleBlock(t *testing.T, chain *core.BlockChain, engine *ethash.Ethash, blockNumber int64, blockHash string, curHeader *types.Header, prevHeader *types.Header, curBody *types.Body, block *types.Block) {
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
		err := engine.VerifyHeader(chain, curHeader, true)
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
	origHeader := block.Header()
	if !equal(curHeader, origHeader) {
		t.Fatalf("Header mismatch: have %+v, want %+v", curHeader, origHeader)
	}
	//origBody := block.Body()
	//if !equal(curBody, origBody) {
	//	t.Fatalf("Body mismatch: have %+v, want %+v", curBody, origBody)
	//}
}

func handleAccount(t *testing.T, blockNumber int64, blockHash string, address string, stateDump *state.DumpAccount, origStateDb *state.StateDB) {
	var addr common.Address = common.BytesToAddress(common.FromHex(address))

	origStateBalance := origStateDb.GetBalance(addr)
	stateBalance := new(big.Int)
	stateBalance, ok := stateBalance.SetString(stateDump.Balance, 10)
	if !ok {
		t.Fatalf("SetString: error, value %s", stateDump.Balance)
	}
	if stateBalance.Cmp(origStateBalance) != 0 {
		t.Fatalf("State balance mismatch: have %x, want %x, address %s", origStateBalance, stateBalance, address)
	}
}

func handleReceipts(t *testing.T, blockNumber int64, receipts types.Receipts, header *types.Header) {
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		t.Fatalf("invalid bloom (remote: %x  local: %x), block_number: %x", header.Bloom, rbloom, blockNumber)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveSha(receipts)
	if receiptSha != header.ReceiptHash {
		t.Fatalf("invalid receipt root hash (remote: %x local: %x), block_number: %x", header.ReceiptHash, receiptSha, blockNumber)
	}
}

func createTestTables(t *testing.T, connectionString string, userName string, dbName string) {
	db, dbErr := sql.Open("postgres", connectionString)
	if dbErr != nil {
		t.Fatalf("Filed to open database. %s", dbErr.Error())
	}
	defer db.Close()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory. %s", err)
	}

	dir = path.Dir(dir + "/../extdb/migrations/")
	//cmd := exec.Command("goose", "-dir", dir, "postgres", "\"user=" + userName + " dbname=" + dbName + " sslmode=disable\"", "up")
	//if err := cmd.Run(); err != nil {
	//	t.Fatalf("Failed to exec goose. %s", err)
	//}
}

func createTestDatabase(t *testing.T, connectionString string, noCreateDb bool, noDropDb bool) (string, func()) {
	var (
		db     *sql.DB
		dbErr  error
		dbName string
	)

	u, err := url.Parse(connectionString)

	if !noCreateDb {
		if err != nil {
			t.Fatalf("Failed to parse connection string. %s", err.Error())
		}

		db, dbErr = sql.Open("postgres", connectionString)
		if dbErr != nil {
			t.Fatalf("Filed to open database. %s", dbErr.Error())
		}

		rand.Seed(time.Now().UnixNano())
		dbName = "jsearch" + strconv.FormatInt(rand.Int63(), 10)

		_, err = db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			t.Fatalf("Failed to create database. %s", err.Error())
		}

		connectionString = u.Scheme + "://" + u.User.String() + "@" + u.Host + "/" + dbName + "?" + u.RawQuery
	}

	createTestTables(t, connectionString, u.User.String(), dbName)

	return connectionString, func() {
		if !noDropDb {
			_, err := db.Exec("DROP DATABASE " + dbName)
			if err != nil {
				t.Fatalf("Drop database failed. %s", err.Error())
			}
		}
		db.Close()
	}
}

func generateBlockchain(t *testing.T) *core.BlockChain {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addrs   = []common.Address{
			crypto.PubkeyToAddress(key1.PublicKey),
			crypto.PubkeyToAddress(key2.PublicKey),
			crypto.PubkeyToAddress(key3.PublicKey),
		}
		testdb = ethdb.NewMemDatabase()
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
			tx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[0]), addrs[1], big.NewInt(100000), params.TxGas*2., big.NewInt(1), []byte{0x11, 0x11, 0x11}), signer, key1)
			gen.SetCoinbase(addrs[2])
			gen.SetExtra([]byte("addr3"))
			gen.AddTx(tx)
		case 1:
			// In block 2, addr1 sends some more ether to addr2.
			// addr2 passes it on to addr3.
			tx1, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[0]), addrs[1], big.NewInt(1000), params.TxGas*3, big.NewInt(1), []byte{0x22, 0x22, 0x22}), signer, key1)
			tx2, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addrs[1]), addrs[2], big.NewInt(1000), params.TxGas*3, big.NewInt(1), []byte{0x33, 0x33, 0x33}), signer, key2)
			gen.SetCoinbase(addrs[0])
			gen.SetExtra([]byte("addr1"))
			gen.AddTx(tx1)
			gen.AddTx(tx2)
			b0 := gen.PrevBlock(0).Header()
			b0.Extra = []byte("uncle1")
			gen.AddUncle(b0)
		case 2:
			// Block 3 is empty but was mined by addr3.
			gen.SetCoinbase(addrs[2])
			gen.SetExtra([]byte("addr3"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			gen.SetCoinbase(addrs[1])
			gen.SetExtra([]byte("addr2"))
			b2 := gen.PrevBlock(1).Header()
			b2.Extra = []byte("uncle2")
			gen.AddUncle(b2)
			b3 := gen.PrevBlock(2).Header()
			b3.Extra = []byte("uncle3")
			gen.AddUncle(b3)
		}
	})

	// Import the chain. This runs all block validation rules.
	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(testdb, nil, gspec.Config, engine, vm.Config{}, nil)

	if i, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("insert error (block %d): %v\n", blocks[i].NumberU64(), err)
	}

	return blockchain
}

// Tests that blockchain saving works.
func TestBlockchainSaving(t *testing.T) {
	var (
		err error
		db  *sql.DB
	)

	connectionString, dropDb := createTestDatabase(t, *extdb_constr, *extdb_nocreatedb, *extdb_nodropdb)
	err = extdb.NewExtDBpg(connectionString)
	if err != nil {
		t.Fatalf("Filed to open database. %s", err.Error())
	}

	blockchain := generateBlockchain(t)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Filed to open database. %s", err.Error())
	}
	readBlocks(t, db, blockchain)
	readAccounts(t, db, blockchain)
	blockchain.Stop()
	db.Close()

	extdb.Close()

	dropDb()
}
