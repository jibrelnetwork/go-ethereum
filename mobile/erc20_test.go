package geth

import (
	"bytes"
	"encoding/hex"
	"testing"
)

/*
func TestSmokeContractABI(t *testing.T) {
	abi, err := NewEIP20ABI()
	check(t, err, "Can't instantiate ContractABI")
	params := NewParameters(2)
	addr, err := NewAddressFromHex("0x0c7c5d1ac0b51a7ec44a52e28c18598e7fcdf32a")
	check(t, err, "Can't create address")
	params.SetAddress(0, addr)
	params.SetBigInt(1, NewBigInt(23))
	data, err := abi.PackArguments("transfer", params)
	check(t, err, "Can't pack params")
	if len(data) == 0 {
		t.Errorf("Packed params is empty")
	}
}

func TestContractCalling(t *testing.T) {
	conn, err := NewEthereumClient("ws://192.168.0.11:8545")
	check(t, err, "can't connect")
	addr, err := NewAddressFromHex("0x87e557f7a46e385b65d5b4b038d3988eb30288fc")
	check(t, err, "can't create addr")
	contract, err := NewContract(addr, conn)
	check(t, err, "can't crate contract")
	name, err := contract.Name(nil)
	check(t, err, "can't get name")
	symbol, err := contract.Symbol(nil)
	check(t, err, "can't get sybmol")
	t.Errorf("%v %v", name, symbol)
}
*/

func TestKeccak256(t *testing.T) {
	refString := "a4007fbb6da9486dd06419a9bc3dda4c7ee16824168d5e51611d370ca14ad906"
	reference, err := bytesFromString(refString)
	check(t, err, "reference string malformed")

	result := Keccak256Hash([]byte("Test hash string"))
	//fmt.Printf("wait:   %v\nreturn: %v\n", refString, hex.EncodeToString(result[:]))
	if !bytes.Equal(result[:], reference) {
		t.Errorf("expected hash: %v\nreceived hash: %v", reference, result)
	}
}

func TestAddressFromPublicKey(t *testing.T) {
	pubKeyString := "03ace171cd7b204b205867076a9776c80502904c0345d00cc811a8a76446b0095502c8b571debb0fd00e653a75bad41834805dbd7ce7277101d34013d52ec70a8f"
	pubKey, err := bytesFromString(pubKeyString)
	check(t, err, "pubkey malformed")
	refString := "FDFabEB69cC54A55F0fC220c94c13E817ed34b4C"
	reference, err := bytesFromString(refString)
	check(t, err, "reference string malformed")

	result := NewAddressFromPublicKey(pubKey)
	if !bytes.Equal(result.GetBytes(), reference) {
		t.Errorf("\nwait:   %v\nreturn: %v\n", hex.EncodeToString(reference), hex.EncodeToString(result.GetBytes()))
	}
}

func bytesFromString(s string) ([]byte, error) {
	referenceString := []byte(s)
	reference := make([]byte, hex.DecodedLen(len(referenceString)))
	_, err := hex.Decode(reference, referenceString)
	return reference, err
}

func check(t *testing.T, err error, description string) {
	if err != nil {
		t.Errorf("error %v: %v", description, err)
	}
}
