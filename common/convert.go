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
			fmt.Println("Error: ", err)
			return ""
		}
		if op, ok := opcodes[int(code)]; ok {
			result = result + op.Name

			if op.Immediates == 1 {
				immediate := bytecode[i+2 : i+2+(op.Immediates*2)]
				immediateInt, err := strconv.ParseInt(immediate, 16, 64)
				immediate = strconv.Itoa(int(immediateInt))

				if err != nil {
					fmt.Println(err)
				}
				result = result + "(" + string(immediate) + ")"
			} else if op.Immediates > 1 {
				fmt.Println("op.Code:", op.Code)
				fmt.Println("op.Name:", op.Name)
				fmt.Println("op.Immediates:", op.Immediates)
				immediate := bytecode[i+2 : i+2+(op.Immediates*2)]
				result = result + "(0x" + immediate + ")"
			}
			i += (op.Immediates * 2)
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
	immediate = strings.ToLower(immediate)

	op := opcodes[opcode]

	if op.Name == "" {
		errmsg := "Opcode not assigned: " + opcode
		return "", errors.New(errmsg)
	}

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
		if op.Name == "RJUMPV" {
			values := strings.Split(immediate, ",")
			count := len(values)
			result += fmt.Sprintf("%02x", count)
			for _, ro := range values {
				relativeOffset, err := strconv.ParseInt(ro, 10, 64)

				if err != nil {
					return "", errors.New("Error parsing RJUMPV relative offset")
				}

				if relativeOffset < 0 {
					relativeOffset = 65535 - (relativeOffset * -1) + 1
				}
				result += fmt.Sprintf("%04x", relativeOffset)
			}

			return result, nil
		}

		// For >push2 a hexadecimal must be received as parameter
		imm_hex := ""
		if op.Code > 0x60 && op.Code <= 0x7f {
			if immediate[0:2] != "0x" {
				return "", errors.New("hexadecimal value expected")
			}
			immediate = immediate[2:]
			totalBytes := op.Code - 0x5f

			for {
				if len(immediate)/2 == totalBytes {
					break
				}

				immediate = "0" + immediate
			}
			imm_hex = immediate
		} else {
			imm, err := strconv.ParseInt(immediate, 10, 64)
			if err != nil {
				return "", err
			}

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
				for {
					if len(imm_hex) < op.Immediates*2 {
						imm_hex = "00" + imm_hex
					} else {
						break
					}
				}
			}
		}

		result = result + imm_hex
	}

	return result, nil
}

func Mnem2Evm(mn string) string {
	tokens := strings.Split(mn, " ")
	result := ""

	for _, token := range tokens {
		token = strings.Trim(token, " ")
		if token == "" {
			continue
		}
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
			return ""
		}
		result = result + evm
	}

	return result
}
