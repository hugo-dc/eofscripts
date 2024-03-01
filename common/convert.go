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
	Position int
	OpCode
	Immediates []Immediate
}

func (op OpCall) ToBytecode() (string, error) {
	bytecode := op.OpCode.AsHex()

	for _, im := range op.Immediates {
		if im.Type == Label {
			return "", errors.New("Not expected label " + im.Immediate)
		}

		bytecode += im.Immediate
	}

	return bytecode, nil
}

type EOFObjectModifier struct {
	Magic       bool
	Version     bool
	TypeHeader  bool
	CodeHeader  bool
	DataHeader  bool
	Terminator  bool
	TypeSection map[int]string
	CodeSection map[int]bool
	DataSection bool
}

func DescribeBytecode(bytecode string) ([]OpCall, error) {
	result := make([]OpCall, 0)
	opcodes := GetOpcodesByNumber()

	for i := 0; i < len(bytecode); i += 2 {
		code_str := bytecode[i:(i + 2)]
		code, err := strconv.ParseInt(code_str, 16, 64)

		if err != nil {
			return nil, err
		}

		if op, ok := opcodes[int(code)]; ok {
			opCall := OpCall{Position: i / 2, OpCode: op, Immediates: make([]Immediate, 0)}

			if op.Name == "" {
				return nil, errors.New(fmt.Sprintf("Opcode not found: %d", op.Code))
			}

			if op.Immediates > 0 {
				immediate := bytecode[i+2 : i+2+(op.Immediates*2)]
				immediateInt := int64(0)
				if op.Immediates < 8 {
					immediateInt_tmp, err := strconv.ParseInt(immediate, 16, 64)

					if err != nil {
						return nil, err
					}
					immediate = fmt.Sprintf("%0*x", op.Immediates*2, immediateInt_tmp)
					immediateInt = immediateInt_tmp
				}
				opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immediate})

				// RJUMPV can have many immediates
				if op.Name == "RJUMPV" {
					i += 2
					for j := 0; j <= int(immediateInt); j++ {
						if i+6 > len(bytecode) {
							return nil, errors.New(fmt.Sprintf("Truncated RJUMPV at position %d", i))
						}
						immediate := bytecode[i+2 : i+6]
						opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immediate})
						i += 4
					}
					i -= 2
				}
			}
			result = append(result, opCall)

			i += op.Immediates * 2
		} else {
			return nil, errors.New(fmt.Sprintf("Opcode not found: %s", code_str))
		}
	}

	return result, nil
}

func Evm2Mnem(opcalls []OpCall) string {
	result := ""
	for _, op := range opcalls {
		result += op.Name

		if op.OpCode.Immediates == 1 {
			immInt, err := strconv.ParseInt(op.Immediates[0].Immediate, 16, 64)

			if err != nil {
				panic(err)
			}

			if op.Name == "RJUMPV" {
				for i := 0; i <= int(immInt); i++ {
					immInt2, err := strconv.ParseInt(op.Immediates[i+1].Immediate, 16, 64)

					if err != nil {
						panic(err)
					}

					if immInt2 > 32767 {
						immInt2 = ((65535 - immInt2) + 1) * -1
					}

					if immInt == 0 {
						result = result + fmt.Sprintf("(%d)", immInt2)
					} else if i == 0 {
						result = result + fmt.Sprintf("(%d,", immInt2)
					} else if i == int(immInt) {
						result = result + fmt.Sprintf("%d)", immInt2)
					} else {
						result = result + fmt.Sprintf("%d,", immInt2)
					}
				}
			} else {
				result = result + fmt.Sprintf("(%d)", immInt)
			}
		} else if op.OpCode.Immediates > 1 {
			immediate := op.Immediates[0].Immediate
			imm, err := strconv.ParseInt(immediate, 16, 64)

			if err != nil {
				result = result + "(0x" + immediate + ")"
			} else {
				if imm > 32767 {
					imm = ((65535 - imm) + 1) * -1
				}
				result = result + "(" + strconv.Itoa(int(imm)) + ")"
			}
		}

		result += " "
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
			max_size := len(values) - 1
			opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: fmt.Sprintf("%02x", max_size)})
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
			opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: imm_hex})
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
	multiplying := false
	for _, token := range tokens {
		token = strings.Trim(token, " ")
		if token == "" {
			continue
		}
		if token[len(token)-1] == ':' {
			labels[token[:len(token)-1]] = pos
			continue
		}

		if token == "*" {
			multiplying = true
			continue
		}

		if multiplying {
			multiplying = false
			value, err := strconv.ParseInt(token, 10, 64)
			if err != nil {
				fmt.Println(err)
				return ""
			}
			if value < 1 {
				fmt.Println("Error, invalid multiplier: ", token)
				return ""
			}
			prevOpCall := evm[len(evm)-1]
			for i := 0; i < int(value)-1; i++ {
				evm = append(evm, prevOpCall)
			}
			pos += int(value) - 1
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
			return ""
		}
		evm = append(evm, opCall)
		if opCall.OpCode.Name == "RJUMPV" {
			pos += (len(opCall.Immediates) - 1) * 2
		}
		pos = pos + 1 + opCall.OpCode.Immediates
	}

	result := ""
	pos = 0
	for _, op := range evm {
		result += op.OpCode.AsHex()

		//fmt.Println(pos, op.OpCode)
		pos = pos + 1 + op.OpCode.Immediates
		if op.OpCode.Name == "RJUMPV" {
			pos += (len(op.Immediates) - 1) * 2 // Add the immediates
		}

		for _, im := range op.Immediates {
			if im.Type == Label {
				if p, ok := labels[im.Immediate]; ok {
					if op.OpCode.Name == "RJUMP" || op.OpCode.Name == "RJUMPI" || op.OpCode.Name == "RJUMPV" {
						//fmt.Println("pos:", pos)
						//fmt.Println("p:", p)
						//fmt.Println("p-pos:", p-pos)
						if p >= pos {
							result += fmt.Sprintf("%04x", p-pos)
						} else {
							imm_hex := strconv.FormatUint(uint64(p-pos), 16)
							if op.OpCode.Name == "RJUMPV" {
								imm_hex = imm_hex[len(imm_hex)-4:]
							} else {
								imm_hex = imm_hex[len(imm_hex)-op.OpCode.Immediates*2:]
							}
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

func NewEOFObjectModifier() EOFObjectModifier {
	return EOFObjectModifier{
		Magic:       true,
		Version:     true,
		TypeHeader:  true,
		CodeHeader:  true,
		DataHeader:  true,
		Terminator:  true,
		TypeSection: make(map[int]string),
		CodeSection: make(map[int]bool),
		DataSection: true,
	}
}

func ModifyEOFObject(eofObject EOFObject, modifier EOFObjectModifier) string {
	newcode := ""
	if modifier.Magic {
		newcode += "ef00"
	}
	if modifier.Version {
		newcode += fmt.Sprintf("%02x", eofObject.Version)
	}
	if modifier.TypeHeader {
		newcode += fmt.Sprintf("01%04x", len(eofObject.Types)*4)
	}
	if modifier.CodeHeader {
		newcode += fmt.Sprintf("02%04x", len(eofObject.CodeSections))

		for _, cs := range eofObject.CodeSections {
			newcode += fmt.Sprintf("%04x", len(cs)/2)
		}
	}
	if modifier.DataHeader {
		newcode += fmt.Sprintf("03%04x", len(eofObject.Data)/2)
	}
	if modifier.Terminator {
		newcode += "00"
	}

	for i, t := range eofObject.Types {
		if s, ok := modifier.TypeSection[i]; ok {
			// If it is blank, means has to be removed
			if s != "" {
				newcode += s
			}
		} else {
			newcode += fmt.Sprintf("%02x%02x%02x", t[0], t[1], t[2])
		}
	}

	for i, cs := range eofObject.CodeSections {
		if remain, ok := modifier.CodeSection[i]; ok {
			if remain {
				newcode += cs
			}
		} else {
			newcode += cs
		}
	}

	if modifier.DataSection {
		newcode += eofObject.Data
	}
	return newcode
}
