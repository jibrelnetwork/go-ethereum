# Mobile bindings for go-ethereum lib

#### Packing params for contract method calling
Use `ContractABI` to pack parameters to byte array for contract method calling. `Parameters` is wrapper for go array of params. Example in `erc20_test.go` `TestSmokeContractABI`.
Example in go:
```
	abi, _ := NewEIP20ABI()
	params := NewParameters(2)
	addr, _ := NewAddressFromHex("0x0c7c5d1ac0b51a7ec44a52e28c18598e7fcdf32a")
	params.SetAddress(0, addr)
	params.SetBigInt(1, NewBigInt(23))
	data, _ := abi.PackArguments("transfer", params)
```
example in swift:
```
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
