package common

import (
	"encoding/hex"
	"testing"
)

var mnems = [][]string{
	{"RJUMPI", "-4", "e1fffc"},
	{"RJUMP", "0", "e00000"},
	{"PUSH1", "0", "6000"},
	{"STOP", "", "00"},
	{"CODESIZE", "", "38"},
	{"RJUMPV", "0,0", "e20200000000"},
	{"RJUMPV", "0,1", "e20200000001"},
	{"RJUMPV", "1,0", "e20200010000"},
	{"RJUMPV", "1,1", "e20200010001"},
	{"RJUMPV", "1,2", "e20200010002"},
	{"RJUMPV", "-1,0", "e202ffff0000"},
	{"RJUMPV", "-1,1", "e202ffff0001"},
	{"RJUMPV", "-1,-1", "e202ffffffff"},
	{"RJUMPV", "0,0,0", "e203000000000000"},
}

var bcodes = []string{
	"5dfffc",
	"5c0000",
	"6001",
	"00",
}

/*
func TestDescribeBytecode(t *testing.T) {
	for _, bc := range bcodes {
		result := DescribeBytecode(bc)

		if result != "x" {
			t.Errorf("Result %v", result)
		}

			r_bc, err := result[0].ToBytecode()

			if err != nil {
				t.Errorf("Not expected error %v", err)
			}

			if r_bc != bc {
				t.Errorf("Bytecode %v (%v) not as expected: %v", r_bc, result, bc)
			}
	}
}
*/

func TestOpcode2EVM(t *testing.T) {
	for _, mn := range mnems {
		result, err := opcode2evm(mn[0], mn[1])
		if err != nil {
			t.Errorf("Not expected error: %v", err)
		}

		result_str := result.OpCode.AsHex()
		for _, im := range result.Immediates {
			if im.Type == Label {
				t.Errorf("Not expected label: %v", im.Immediate)
			}
			result_str += hex.EncodeToString(im.Immediate)
		}

		if result_str != mn[2] {
			t.Errorf("EVM Code %v not as expected: %v", result_str, mn[2])
		}
	}
}
