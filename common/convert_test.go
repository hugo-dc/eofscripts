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
	{"RJUMPV", "0200000000", "5e0200000000"},
}

func TestOpcode2EVM(t *testing.T) {
	for _, mn := range mnems {
		result, err := opcode2evm(mn[0], mn[1])

		if err != nil {
			t.Errorf("Not expected error %v", err)
		}

		if result != mn[2] {
			t.Errorf("EVM Code %v not as expected: %v", result, mn[2])
		}
	}
}
