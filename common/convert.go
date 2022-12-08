package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Evm2Mnem(bytecode string) string {
	result := ""
	opcodes := GetOpcodesByNumber()
	for i := 0; i < len(bytecode); i += 2 {
		code_str := bytecode[i:(i + 2)]
		code, err := strconv.ParseInt(code_str, 16, 64)

		if err != nil {
			// TODO
		}
		if op, ok := opcodes[int(code)]; ok {
			result = result + op.Name

			if op.Immediates > 0 {
				if op.Name != "RJUMP" && op.Name != "RJUMPI" {
					immediate := bytecode[i+2 : i+2+(op.Immediates*2)]
					result = result + "(0x" + immediate + ")"
					i += (op.Immediates * 2)
				}
			}
		} else {
			fmt.Println("Error: opcode " + "0x" + code_str + " not found")
			return ""
		}
		result = result + " "
	}
	return result
}

func opcode2evm(opcode string, immediate string) (string, error) {
	opcodes := GetOpcodesByName()

	op := opcodes[opcode]
	if immediate == "" && op.Immediates > 0 {
		return "", errors.New("Immediate expected for opcode")
	}

	if immediate != "" && op.Immediates == 0 {
		return "", errors.New("Immediate NOT expected for opcode")
	}

	result := strconv.FormatInt(int64(op.Code), 16)
	if len(result)%2 != 0 {
		result = "0" + result
	}

	if immediate != "" {
		imm, err := strconv.ParseInt(immediate, 10, 64)

		if err != nil {
			return "", err
		}

		imm_hex := ""
		if imm < 0 {
			if op.Name != "RJUMP" && op.Name != "RJUMPI" {
				return "", errors.New("Negative immediate only possible for RJUMP and RJUMPI")
			}

			imm_hex = strconv.FormatUint(uint64(imm), 16)
			imm_hex = imm_hex[len(imm_hex)-op.Immediates*2:]
		} else {
			imm_hex = strconv.FormatInt(imm, 16)
			if len(imm_hex)%2 != 0 {
				imm_hex = "0" + imm_hex
			}
		}

		result = result + imm_hex
	}

	return result, nil
}

func Mnem2Evm(mn string) string {
	tokens := strings.Split(mn, " ")
	result := ""

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		opcode := ""
		immediate := ""
		if strings.Contains(token, "(") {
			elements := strings.Split(token, "(")
			opcode = elements[0]
			immediate = elements[1][0 : len(elements[1])-1]
		} else {
			opcode = token
		}

		evm, err := opcode2evm(opcode, immediate)

		if err != nil {
			fmt.Println("Error: ", err, opcode, immediate)
		}

		result = result + evm
	}

	return result
}
