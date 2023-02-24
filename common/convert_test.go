package common

import (
	"testing"
)

var mnems = [][]string{
	{"RJUMPI", "-4", "5dfffc"},
	{"RJUMP", "0", "5c0000"},
	{"PUSH1", "0", "6000"},
	{"STOP", "", "00"},
	{"CODESIZE", "", "38"},
	{"RJUMPV", "0,0", "5e0200000000"},
}

var bcodes = []string{
	"5dfffc",
	"5c0000",
	"6001",
	"00",
}

func TestDescribeBytecode(t *testing.T) {
	for _, bc := range bcodes {
		result, err := DescribeBytecode(bc)

		if err != nil {
			t.Errorf("Not expected error %v", err)
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

func TestOpcode2EVM(t *testing.T) {
	for _, mn := range mnems {
		result, err := opcode2evm(mn[0], mn[1])

		if err != nil {
			t.Errorf("Not expected error %v", err)
		}

		result_str := result.OpCode.AsHex()

		for _, im := range result.Immediates {
			if im.Type == Label {
				t.Errorf("Not expected label %v", im.Immediate)
			}
			result_str += im.Immediate
		}

		if result_str != mn[2] {
			t.Errorf("EVM Code %v not as expected: %v", result_str, mn[2])
		}
	}
}
