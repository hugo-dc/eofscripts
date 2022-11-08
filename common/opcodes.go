package common

import "strconv"

type OpCode struct {
	Code   int
	HexStr string
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
