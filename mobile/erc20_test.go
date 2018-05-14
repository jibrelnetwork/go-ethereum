package geth

import (
	"testing"
)

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

func check(t *testing.T, err error, description string) {
	if err != nil {
		t.Errorf("error %v: %v", description, err)
	}
}
