package common

import "strconv"

type OpCode struct {
	Name       string
	Code       int
	Immediates int
}

func Stop() OpCode {
	return OpCode{
		Code: 0x00,
	}
}

func CodeCopy() OpCode {
	return OpCode{
		Code: 0x39,
	}
}

func MStore() OpCode {
	return OpCode{
		Code: 0x52,
	}
}

func Push1() OpCode {
	return OpCode{
		Code: 0x60,
	}
}

func Return() OpCode {
	return OpCode{
		Code: 0xf3,
	}
}

func Shl() OpCode {
	return OpCode{
		Code: 0x1b,
	}
}

func (op OpCode) AsHex() string {
	opValue := strconv.FormatInt(int64(op.Code), 16)

	if len(opValue)%2 != 0 {
		opValue = "0" + opValue
	}

	return opValue
}

var opcodes = []OpCode{
	{Name: "STOP", Code: 0x00, Immediates: 0},
	{Name: "ADD", Code: 0x01, Immediates: 0},
	{Name: "MUL", Code: 0x02, Immediates: 0},
	{Name: "SUB", Code: 0x03, Immediates: 0},
	{Name: "DIV", Code: 0x04, Immediates: 0},
	{Name: "SDIV", Code: 0x05, Immediates: 0},
	{Name: "MOD", Code: 0x06, Immediates: 0},
	{Name: "SMOD", Code: 0x07, Immediates: 0},
	{Name: "LT", Code: 0x10, Immediates: 0},
	{Name: "GT", Code: 0x11, Immediates: 0},
	{Name: "SLT", Code: 0x12, Immediates: 0},
	{Name: "EQ", Code: 0x14, Immediates: 0},
	{Name: "ISZERO", Code: 0x15, Immediates: 0},
	{Name: "AND", Code: 0x16, Immediates: 0},
	{Name: "SHL", Code: 0x1b, Immediates: 0},
	{Name: "SHR", Code: 0x1c, Immediates: 0},
	{Name: "SAR", Code: 0x1d, Immediates: 0},
	{Name: "ADDRESS", Code: 0x30, Immediates: 0},
	{Name: "CALLDATALOAD", Code: 0x35, Immediates: 0},
	{Name: "CODECOPY", Code: 0x39, Immediates: 0},
	{Name: "POP", Code: 0x50, Immediates: 0},
	{Name: "MSTORE", Code: 0x52, Immediates: 0},
	{Name: "MSTORE8", Code: 0x53, Immediates: 0},
	{Name: "SSTORE", Code: 0x55, Immediates: 0},
	{Name: "JUMP", Code: 0x56, Immediates: 0},
	{Name: "JUMPI", Code: 0x5b, Immediates: 0},
	{Name: "RJUMP", Code: 0x5c, Immediates: 2},
	{Name: "RJUMPI", Code: 0x5d, Immediates: 2},
	{Name: "PUSH1", Code: 0x60, Immediates: 1},
	{Name: "PUSH2", Code: 0x61, Immediates: 2},
	{Name: "PUSH3", Code: 0x62, Immediates: 3},
	{Name: "PUSH4", Code: 0x63, Immediates: 4},
	{Name: "PUSH5", Code: 0x64, Immediates: 5},
	{Name: "PUSH6", Code: 0x65, Immediates: 6},
	{Name: "PUSH7", Code: 0x66, Immediates: 7},
	{Name: "PUSH8", Code: 0x67, Immediates: 8},
	{Name: "PUSH9", Code: 0x68, Immediates: 9},
	{Name: "PUSH10", Code: 0x69, Immediates: 10},
	{Name: "PUSH11", Code: 0x6a, Immediates: 11},
	{Name: "PUSH12", Code: 0x6b, Immediates: 12},
	{Name: "PUSH13", Code: 0x6c, Immediates: 13},
	{Name: "PUSH14", Code: 0x6d, Immediates: 14},
	{Name: "PUSH15", Code: 0x6e, Immediates: 15},
	{Name: "PUSH16", Code: 0x6f, Immediates: 16},
	{Name: "PUSH17", Code: 0x70, Immediates: 17},
	{Name: "PUSH18", Code: 0x71, Immediates: 18},
	{Name: "PUSH19", Code: 0x72, Immediates: 19},
	{Name: "PUSH20", Code: 0x73, Immediates: 20},
	{Name: "PUSH21", Code: 0x74, Immediates: 21},
	{Name: "PUSH22", Code: 0x75, Immediates: 22},
	{Name: "PUSH23", Code: 0x76, Immediates: 23},
	{Name: "PUSH24", Code: 0x77, Immediates: 24},
	{Name: "PUSH25", Code: 0x78, Immediates: 25},
	{Name: "PUSH26", Code: 0x79, Immediates: 26},
	{Name: "PUSH27", Code: 0x7a, Immediates: 27},
	{Name: "PUSH28", Code: 0x7b, Immediates: 28},
	{Name: "PUSH29", Code: 0x7c, Immediates: 29},
	{Name: "PUSH30", Code: 0x7d, Immediates: 30},
	{Name: "PUSH31", Code: 0x7e, Immediates: 31},
	{Name: "PUSH32", Code: 0x7f, Immediates: 32},
	{Name: "DUP1", Code: 0x80, Immediates: 0},
	{Name: "DUP2", Code: 0x81, Immediates: 0},
	{Name: "DUP3", Code: 0x82, Immediates: 0},
	{Name: "DUP4", Code: 0x83, Immediates: 0},
	{Name: "DUP5", Code: 0x84, Immediates: 0},
	{Name: "SWAP1", Code: 0x90, Immediates: 0},
	{Name: "SWAP2", Code: 0x91, Immediates: 0},
	{Name: "SWAP3", Code: 0x92, Immediates: 0},
	{Name: "CALLF", Code: 0xb0, Immediates: 2},
	{Name: "RETF", Code: 0xb1, Immediates: 0},
	{Name: "JUMPF", Code: 0xb2, Immediates: 2},
	{Name: "RETURN", Code: 0xf3, Immediates: 0},
	{Name: "RETF", Code: 0xfc, Immediates: 0},
	{Name: "REVERT", Code: 0xfd, Immediates: 0},
	{Name: "INVALID", Code: 0xfe, Immediates: 0},
	{Name: "SELFDESTRUCT", Code: 0xff, Immediates: 0},
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
	for i := 0; i < len(opcodes); i++ {
		result[opcodes[i].Code] = opcodes[i]
	}
	return result
}
