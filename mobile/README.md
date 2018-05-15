# Mobile bindings for go-ethereum lib

#### Packing params for contract method calling
Use `ContractABI` to pack parameters to byte array for contract method calling. `Parameters` is wrapper for go array of params. Example in `erc20_test.go` `TestSmokeContractABI`.
Example in go:
```go
abi, _ := NewEIP20ABI()
params := NewParameters(2)
addr, _ := NewAddressFromHex("0x0c7c5d1ac0b51a7ec44a52e28c18598e7fcdf32a")
params.SetAddress(0, addr)
params.SetBigInt(1, NewBigInt(23))
data, _ := abi.PackArguments("transfer", params)
```
example in swift:
```swift
do {
	 let params = GethNewParameters(2)!
	 var error: NSError? = nil
	 let wallet1 = GethNewAddressFromHex("0x4eb557a7875019332d4fedd4fb8fc209e455b328", &error)!
	 try params.setAddress(0, address: wallet1)
	 try params.setBigInt(1, bigInt: GethNewBigInt(100))

	 let abi = GethNewEIP20ABI(&error)
	 let data = try abi?.packArguments("transfer", params: params)
	 print("data: \(data)")
} catch {
}

```

#### Calling ERC20-compatible contract methods
Use `Contract` to call methods and send transactions. You need connection to node to make calls and sending transactions:
```swift
let ctx = GethNewContext()!
let client = GethNewEthereumClient("ws://192.168.0.11:8545", &error)!
let chainId = try client.networkID(ctx) // for EIP155 support
```
for call methods
```swift
var balance = try contract.balance(of: wallet1, opts: nil)
```
for send transactions
```swift
let gasPrice = try client.suggestGasPrice(ctx)
let noncePtr = UnsafeMutablePointer<Int64>.allocate(capacity: 1)
try client.getPendingNonce(at: ctx, account: wallet1, nonce: noncePtr)
let acc1 = account(address: wallet1.getHex())
try keystore.unlock(acc1, passphrase: "test")
let topts = GethNewTransactOpts(ctx, self, wallet1, GethNewBigInt(noncePtr.pointee), nil/*GethNewBigInt(16) for EIP155*/, gasPrice, GethNewBigInt(3000000)) // self should implement `GethSignerProtocol`

let transaction = try contract.transfer(topts, to: wallet2, value: GethNewBigInt(3000))
```
exmaple of `GethSignerProtocol` implementation
```swift
extension ViewController: GethSignerProtocol {
  func sign(_ p0: GethAddress!, p1: GethTransaction!) throws -> GethTransaction {
    if let account = account(address: p0.getHex()) {
      do {
        return try keystore.signTx(account, tx: p1, chainID: nil)
      } catch {
        print("can't sign \(p1), error: \(error)")
      }
    }
    return p1
  }
}
```