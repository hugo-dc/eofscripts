package common

import (
	"strconv"
)

type OpCode struct {
	Name          string
	Code          int
	Immediates    int
	StackInput    int
	StackOutput   int
	IsTerminating bool
}

func (op OpCode) AsHex() string {
	opValue := strconv.FormatInt(int64(op.Code), 16)

	if len(opValue)%2 != 0 {
		opValue = "0" + opValue
	}

	return opValue
}

func (op OpCode) IsPush() bool {
	return op.Code >= 0x60 && op.Code <= 0x7f
}

func (op OpCode) IsEOFDeprecated() bool {
	if op.Name == "JUMP" || op.Name == "JUMPI" || op.Name == "PC" || op.Name == "CALLCODE" || op.Name == "SELFDESTRUCT" {
		return true
	}
	return false
}

var opcodes = []OpCode{
	{Name: "STOP", Code: 0x00, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: true},
	{Name: "ADD", Code: 0x01, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "MUL", Code: 0x02, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SUB", Code: 0x03, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "DIV", Code: 0x04, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SDIV", Code: 0x05, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "MOD", Code: 0x06, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SMOD", Code: 0x07, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "ADDMOD", Code: 0x08, Immediates: 0, StackInput: 3, StackOutput: 1, IsTerminating: false},
	{Name: "MULMOD", Code: 0x09, Immediates: 0, StackInput: 3, StackOutput: 1, IsTerminating: false},
	{Name: "EXP", Code: 0x0a, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SIGNEXTEND", Code: 0x0b, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "LT", Code: 0x10, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "GT", Code: 0x11, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SLT", Code: 0x12, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SGT", Code: 0x13, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "EQ", Code: 0x14, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "ISZERO", Code: 0x15, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "AND", Code: 0x16, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "OR", Code: 0x17, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "XOR", Code: 0x18, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "NOT", Code: 0x19, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "BYTE", Code: 0x1a, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SHL", Code: 0x1b, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SHR", Code: 0x1c, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SAR", Code: 0x1d, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "KECCAK256", Code: 0x20, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "SHA3", Code: 0x20, Immediates: 0, StackInput: 2, StackOutput: 1, IsTerminating: false},
	{Name: "ADDRESS", Code: 0x30, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "BALANCE", Code: 0x31, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "ORIGIN", Code: 0x32, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CALLER", Code: 0x33, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CALLVALUE", Code: 0x34, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CALLDATALOAD", Code: 0x35, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "CALLDATASIZE", Code: 0x36, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CALLDATACOPY", Code: 0x37, Immediates: 0, StackInput: 3, StackOutput: 0, IsTerminating: false},
	{Name: "CODESIZE", Code: 0x38, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CODECOPY", Code: 0x39, Immediates: 0, StackInput: 3, StackOutput: 0, IsTerminating: false},
	{Name: "GASPRICE", Code: 0x3a, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "EXTCODESIZE", Code: 0x3b, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "EXTCODECOPY", Code: 0x3c, Immediates: 0, StackInput: 4, StackOutput: 0, IsTerminating: false},
	{Name: "RETURNDATASIZE", Code: 0x3d, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "RETURNDATACOPY", Code: 0x3e, Immediates: 0, StackInput: 3, StackOutput: 0, IsTerminating: false},
	{Name: "EXTCODEHASH", Code: 0x3f, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "BLOCKHASH", Code: 0x40, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "COINBASE", Code: 0x41, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "TIMESTAMP", Code: 0x42, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "NUMBER", Code: 0x43, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "DIFFICULTY", Code: 0x44, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "GASLIMIT", Code: 0x45, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "CHAINID", Code: 0x46, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "SELFBALANCE", Code: 0x47, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "BASEFEE", Code: 0x48, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "POP", Code: 0x50, Immediates: 0, StackInput: 1, StackOutput: 0, IsTerminating: false},
	{Name: "MLOAD", Code: 0x51, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "MSTORE", Code: 0x52, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: false},
	{Name: "MSTORE8", Code: 0x53, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: false},
	{Name: "SLOAD", Code: 0x54, Immediates: 0, StackInput: 1, StackOutput: 1, IsTerminating: false},
	{Name: "SSTORE", Code: 0x55, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: false},
	{Name: "JUMP", Code: 0x56, Immediates: 0, StackInput: 1, StackOutput: 0, IsTerminating: false},
	{Name: "JUMPI", Code: 0x57, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: false},
	{Name: "PC", Code: 0x58, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "MSIZE", Code: 0x59, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "GAS", Code: 0x5a, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "JUMPDEST", Code: 0x5b, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: false},
	{Name: "NOP", Code: 0x5b, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: false},
	{Name: "RJUMP", Code: 0x5c, Immediates: 2, StackInput: 0, StackOutput: 0, IsTerminating: false},
	{Name: "RJUMPI", Code: 0x5d, Immediates: 2, StackInput: 1, StackOutput: 0, IsTerminating: false},
	{Name: "RJUMPV", Code: 0x5e, Immediates: 1, StackInput: 1, StackOutput: 0, IsTerminating: false},
	{Name: "PUSH0", Code: 0x5f, Immediates: 0, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH1", Code: 0x60, Immediates: 1, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH2", Code: 0x61, Immediates: 2, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH3", Code: 0x62, Immediates: 3, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH4", Code: 0x63, Immediates: 4, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH5", Code: 0x64, Immediates: 5, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH6", Code: 0x65, Immediates: 6, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH7", Code: 0x66, Immediates: 7, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH8", Code: 0x67, Immediates: 8, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH9", Code: 0x68, Immediates: 9, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH10", Code: 0x69, Immediates: 10, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH11", Code: 0x6a, Immediates: 11, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH12", Code: 0x6b, Immediates: 12, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH13", Code: 0x6c, Immediates: 13, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH14", Code: 0x6d, Immediates: 14, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH15", Code: 0x6e, Immediates: 15, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH16", Code: 0x6f, Immediates: 16, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH17", Code: 0x70, Immediates: 17, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH18", Code: 0x71, Immediates: 18, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH19", Code: 0x72, Immediates: 19, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH20", Code: 0x73, Immediates: 20, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH21", Code: 0x74, Immediates: 21, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH22", Code: 0x75, Immediates: 22, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH23", Code: 0x76, Immediates: 23, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH24", Code: 0x77, Immediates: 24, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH25", Code: 0x78, Immediates: 25, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH26", Code: 0x79, Immediates: 26, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH27", Code: 0x7a, Immediates: 27, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH28", Code: 0x7b, Immediates: 28, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH29", Code: 0x7c, Immediates: 29, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH30", Code: 0x7d, Immediates: 30, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH31", Code: 0x7e, Immediates: 31, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "PUSH32", Code: 0x7f, Immediates: 32, StackInput: 0, StackOutput: 1, IsTerminating: false},
	{Name: "DUP1", Code: 0x80, Immediates: 0, StackInput: 1, StackOutput: 2, IsTerminating: false},
	{Name: "DUP2", Code: 0x81, Immediates: 0, StackInput: 2, StackOutput: 3, IsTerminating: false},
	{Name: "DUP3", Code: 0x82, Immediates: 0, StackInput: 3, StackOutput: 4, IsTerminating: false},
	{Name: "DUP4", Code: 0x83, Immediates: 0, StackInput: 4, StackOutput: 5, IsTerminating: false},
	{Name: "DUP5", Code: 0x84, Immediates: 0, StackInput: 5, StackOutput: 6, IsTerminating: false},
	{Name: "DUP6", Code: 0x85, Immediates: 0, StackInput: 6, StackOutput: 7, IsTerminating: false},
	{Name: "DUP7", Code: 0x86, Immediates: 0, StackInput: 7, StackOutput: 8, IsTerminating: false},
	{Name: "DUP8", Code: 0x87, Immediates: 0, StackInput: 8, StackOutput: 9, IsTerminating: false},
	{Name: "DUP9", Code: 0x88, Immediates: 0, StackInput: 9, StackOutput: 10, IsTerminating: false},
	{Name: "DUP10", Code: 0x89, Immediates: 0, StackInput: 10, StackOutput: 11, IsTerminating: false},
	{Name: "DUP11", Code: 0x8a, Immediates: 0, StackInput: 11, StackOutput: 12, IsTerminating: false},
	{Name: "DUP12", Code: 0x8b, Immediates: 0, StackInput: 12, StackOutput: 13, IsTerminating: false},
	{Name: "DUP13", Code: 0x8c, Immediates: 0, StackInput: 13, StackOutput: 14, IsTerminating: false},
	{Name: "DUP14", Code: 0x8d, Immediates: 0, StackInput: 14, StackOutput: 15, IsTerminating: false},
	{Name: "DUP15", Code: 0x8e, Immediates: 0, StackInput: 15, StackOutput: 16, IsTerminating: false},
	{Name: "DUP16", Code: 0x8f, Immediates: 0, StackInput: 16, StackOutput: 17, IsTerminating: false},
	{Name: "SWAP1", Code: 0x90, Immediates: 0, StackInput: 2, StackOutput: 2, IsTerminating: false},
	{Name: "SWAP2", Code: 0x91, Immediates: 0, StackInput: 3, StackOutput: 3, IsTerminating: false},
	{Name: "SWAP3", Code: 0x92, Immediates: 0, StackInput: 4, StackOutput: 4, IsTerminating: false},
	{Name: "SWAP4", Code: 0x93, Immediates: 0, StackInput: 5, StackOutput: 5, IsTerminating: false},
	{Name: "SWAP5", Code: 0x94, Immediates: 0, StackInput: 6, StackOutput: 6, IsTerminating: false},
	{Name: "SWAP6", Code: 0x95, Immediates: 0, StackInput: 7, StackOutput: 7, IsTerminating: false},
	{Name: "SWAP7", Code: 0x96, Immediates: 0, StackInput: 8, StackOutput: 8, IsTerminating: false},
	{Name: "SWAP8", Code: 0x97, Immediates: 0, StackInput: 9, StackOutput: 9, IsTerminating: false},
	{Name: "SWAP9", Code: 0x98, Immediates: 0, StackInput: 10, StackOutput: 10, IsTerminating: false},
	{Name: "SWAP10", Code: 0x99, Immediates: 0, StackInput: 11, StackOutput: 11, IsTerminating: false},
	{Name: "SWAP11", Code: 0x9a, Immediates: 0, StackInput: 12, StackOutput: 12, IsTerminating: false},
	{Name: "SWAP12", Code: 0x9b, Immediates: 0, StackInput: 13, StackOutput: 13, IsTerminating: false},
	{Name: "SWAP13", Code: 0x9c, Immediates: 0, StackInput: 14, StackOutput: 14, IsTerminating: false},
	{Name: "SWAP14", Code: 0x9d, Immediates: 0, StackInput: 15, StackOutput: 15, IsTerminating: false},
	{Name: "SWAP15", Code: 0x9e, Immediates: 0, StackInput: 16, StackOutput: 16, IsTerminating: false},
	{Name: "SWAP16", Code: 0x9f, Immediates: 0, StackInput: 17, StackOutput: 17, IsTerminating: false},
	{Name: "LOG0", Code: 0xa0, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: false},
	{Name: "LOG1", Code: 0xa1, Immediates: 0, StackInput: 3, StackOutput: 0, IsTerminating: false},
	{Name: "LOG2", Code: 0xa2, Immediates: 0, StackInput: 4, StackOutput: 0, IsTerminating: false},
	{Name: "LOG3", Code: 0xa3, Immediates: 0, StackInput: 5, StackOutput: 0, IsTerminating: false},
	{Name: "LOG4", Code: 0xa4, Immediates: 0, StackInput: 6, StackOutput: 0, IsTerminating: false},
	{Name: "CALLF", Code: 0xb0, Immediates: 2, StackInput: 0, StackOutput: 0, IsTerminating: false},
	{Name: "RETF", Code: 0xb1, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: true},
	{Name: "JUMPF", Code: 0xb2, Immediates: 2, StackInput: 0, StackOutput: 0, IsTerminating: false},
	{Name: "CREATE", Code: 0xf0, Immediates: 0, StackInput: 3, StackOutput: 1, IsTerminating: false},
	{Name: "CALL", Code: 0xf1, Immediates: 0, StackInput: 7, StackOutput: 1, IsTerminating: false},
	{Name: "CALLCODE", Code: 0xf2, Immediates: 0, StackInput: 7, StackOutput: 1, IsTerminating: false},
	{Name: "RETURN", Code: 0xf3, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: true},
	{Name: "DELEGATECALL", Code: 0xf4, Immediates: 0, StackInput: 6, StackOutput: 1, IsTerminating: false},
	{Name: "CREATE2", Code: 0xf5, Immediates: 0, StackInput: 4, StackOutput: 1, IsTerminating: false},
	{Name: "STATICCALL", Code: 0xfa, Immediates: 0, StackInput: 6, StackOutput: 1, IsTerminating: false},
	{Name: "REVERT", Code: 0xfd, Immediates: 0, StackInput: 2, StackOutput: 0, IsTerminating: true},
	{Name: "INVALID", Code: 0xfe, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: true},
	{Name: "SELFDESTRUCT", Code: 0xff, Immediates: 0, StackInput: 1, StackOutput: 0, IsTerminating: true},
}

func GetOpcodesByName() map[string]OpCode {
	result := map[string]OpCode{}
	for i := 0; i < len(opcodes); i++ {
		result[opcodes[i].Name] = opcodes[i]
	}

	return result
}

func GetOpcodesByNumber() map[int]OpCode {
	result := map[int]OpCode{}
	for _, op := range opcodes {
		result[op.Code] = op
	}
	return result
}
