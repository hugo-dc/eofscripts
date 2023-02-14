package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	Value int = iota
	Label
)

type Immediate struct {
	Type      int
	Immediate string
}

type OpCall struct {
	OpCode
	Immediates []Immediate
}

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

func opcode2evm(opcode string, immediate string) (OpCall, error) {
	opcodes := GetOpcodesByName()
	immediate = strings.ToLower(immediate)

	op := opcodes[opcode]
	opCall := OpCall{OpCode: op, Immediates: make([]Immediate, 0)}

	if op.Name == "" {
		errmsg := "Opcode not assigned: " + opcode
		return opCall, errors.New(errmsg)
	}

	if immediate == "" && op.Immediates > 0 {
		return opCall, errors.New("Immediate expected for opcode")
	}

	if immediate != "" && op.Immediates == 0 {
		return opCall, errors.New("Immediate NOT expected for opcode")
	}

	if immediate != "" {
		if op.Name == "RJUMPV" {
			values := strings.Split(immediate, ",")
			count := len(values)
			opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: fmt.Sprintf("%02x", count)})
			for _, ro := range values {
				relativeOffset, err := strconv.ParseInt(ro, 10, 64)

				if err != nil {
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Label, Immediate: ro})
				} else {

					if relativeOffset < 0 {
						relativeOffset = 65535 - (relativeOffset * -1) + 1
					}
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: fmt.Sprintf("%04x", relativeOffset)})
				}
			}

			return opCall, nil
		}

		// For >push2 a hexadecimal must be received as parameter
		imm_hex := ""
		if op.Code > 0x60 && op.Code <= 0x7f {
			if immediate[0:2] != "0x" {
				return opCall, errors.New("hexadecimal value expected")
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
				opCall.Immediates = append(opCall.Immediates, Immediate{Type: Label, Immediate: immediate})
			} else {
				if imm < 0 {
					if op.Name != "RJUMP" && op.Name != "RJUMPI" {
						return opCall, errors.New("Negative immediate only possible for RJUMP and RJUMPI")
					}

					imm_hex = strconv.FormatUint(uint64(imm), 16)
					imm_hex = imm_hex[len(imm_hex)-op.Immediates*2:]
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: imm_hex})
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
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: imm_hex})
				}
			}
		}
	}
	return opCall, nil
}

func Mnem2Evm(mn string) string {
	tokens := strings.Split(mn, " ")

	labels := make(map[string]int)
	evm := make([]OpCall, 0)
	pos := 0
	for _, token := range tokens {
		token = strings.Trim(token, " ")
		if token == "" {
			continue
		}
		if token[len(token)-1] == ':' {
			labels[token[:len(token)-1]] = pos
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
		opCall, err := opcode2evm(opcode, immediate)

		if err != nil {
			fmt.Println("Error: ", err, opcode, immediate)
			return ""
		}
		evm = append(evm, opCall)

		pos = pos + 1 + opCall.OpCode.Immediates
	}

	result := ""
	pos = 0
	for _, op := range evm {
		result += op.OpCode.AsHex()

		pos = pos + 1 + op.OpCode.Immediates
		for _, im := range op.Immediates {
			if im.Type == Label {
				if p, ok := labels[im.Immediate]; ok {
					if op.OpCode.Name == "RJUMP" || op.OpCode.Name == "RJUMPI" || op.OpCode.Name == "RJUMPV" {
						if p > pos {
							result += fmt.Sprintf("%04x", p-pos)
						} else {
							imm_hex := strconv.FormatUint(uint64(p-pos), 16)
							imm_hex = imm_hex[len(imm_hex)-op.OpCode.Immediates*2:]
							result += imm_hex
						}
					} else {
						result += fmt.Sprintf("%04x", p)
					}
				} else {
					fmt.Println("Error: Label not found: ", im)
					fmt.Println("labels: ", labels)
				}
			} else {
				result += im.Immediate
			}
		}
	}

	return result
}
