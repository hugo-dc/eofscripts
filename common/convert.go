package common

import (
	"encoding/binary"
	"encoding/hex"
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
	Immediate []byte
	Label     string
}

type OpCall struct {
	Position int
	OpCode
	Immediates []Immediate
}

func (im Immediate) EqualTo(im2 Immediate) bool {
	if im.Type != im2.Type {
		return false
	}
	if len(im.Immediate) != len(im2.Immediate) {
		return false
	}
	for i, imm := range im.Immediate {
		if imm != im2.Immediate[i] {
			return false
		}
	}
	return true
}

func (op OpCall) ToBytecode() ([]byte, error) {
	bytecode := make([]byte, 0)
	bytecode = append(bytecode, byte(op.Code))

	for _, im := range op.Immediates {
		if im.Type == Label {
			return bytecode, errors.New("Not expected label " + im.Label)
		}

		bytecode = append(bytecode, im.Immediate...)
	}
	return bytecode, nil
}

func (op OpCall) EqualTo(op2 OpCall) bool {
	if op.OpCode != op2.OpCode {
		return false
	}
	if len(op.Immediates) != len(op2.Immediates) {
		return false
	}
	for i, imm := range op.Immediates {
		if !imm.EqualTo(op2.Immediates[i]) {
			return false
		}
	}
	return true
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

func DescribeBytecode(bytecode string) string {
	result := ""
	opcodes := GetOpcodesByNumber()

	for i := 0; i < len(bytecode); i += 2 {
		code_str := bytecode[i:(i + 2)]
		code, err := strconv.ParseInt(code_str, 16, 64)

		if err != nil {
			result += code_str + "# Unknown opcode"
		}

		if op, ok := opcodes[int(code)]; ok {
			imm := bytecode[i+2 : i+2+(op.Immediates*2)]
			if op.Name == "" {
				result += code_str + "# Unknown opcode"
			} else {
				if op.Immediates > 0 {
					result += code_str + imm + " # " + op.Name + "\n"
				} else {
					result += code_str + " # " + op.Name + "\n"
				}
			}

			i += op.Immediates * 2
		} else {
			uopName := fmt.Sprintf("OPCODE_%02X", code)
			result += uopName + " "
		}
	}
	return result
}

func BytecodeToOpCalls(bytecode []byte) ([]OpCall, error) {
	result := make([]OpCall, 0)
	opcodes := GetOpcodesByNumber()

	for i := 0; i < len(bytecode); i++ {
		code := int(bytecode[i])
		if op, ok := opcodes[code]; ok {
			opCall := OpCall{Position: i, OpCode: op, Immediates: make([]Immediate, 0)}

			if op.Name == "" {
				return nil, errors.New(fmt.Sprintf("Opcode not found: %d", op.Code))
			}

			if op.Immediates > 0 {
				immediate := bytecode[i+1 : i+1+(op.Immediates)]
				opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immediate})
				immediateInt := 0

				// RJUMPV can have many immediates
				if op.Name == "RJUMPV" {
					immediateInt = int(immediate[0])
					i += 1
					for j := 0; j <= int(immediateInt); j++ {
						if i+3 > len(bytecode) {
							return nil, errors.New(fmt.Sprintf("Truncated RJUMPV at position %d", i))
						}
						immediate := bytecode[i+1 : i+3]
						opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immediate})
						i += 2
					}
					i -= 1
				}
			}
			result = append(result, opCall)
			i += op.Immediates
		} else {
			uopName := fmt.Sprintf("OPCODE_%02X", code)
			undefinedOp := OpCode{Code: int(code), Name: uopName, Immediates: 0, StackInput: 0, StackOutput: 0, IsTerminating: false}
			opCall := OpCall{Position: i / 2, OpCode: undefinedOp, Immediates: make([]Immediate, 0)}
			result = append(result, opCall)
		}
	}
	return result, nil
}

func Evm2Asm(opcalls []OpCall, prefix string, imm_indicators []string) string {
	result := ""
	par_op := imm_indicators[0]
	par_cl := imm_indicators[1]
	for _, op := range opcalls {
		result += prefix + op.Name

		if op.OpCode.Immediates == 1 {
			immInt := int(op.Immediates[0].Immediate[0])

			if op.Name == "RJUMPV" {
				for i := 0; i <= int(immInt); i++ {
					immInt2 := int64(binary.BigEndian.Uint16(op.Immediates[i+1].Immediate))

					if immInt2 > 32767 {
						immInt2 = ((65535 - immInt2) + 1) * -1
					}

					if immInt == 0 {
						result = result + fmt.Sprintf("%v%d%v", par_op, immInt2, par_cl)
					} else if i == 0 {
						result = result + fmt.Sprintf("%v%d,", par_op, immInt2)
					} else if i == int(immInt) {
						result = result + fmt.Sprintf("%d%v", immInt2, par_cl)
					} else {
						result = result + fmt.Sprintf("%d,", immInt2)
					}
				}
			} else {
				result = result + fmt.Sprintf("%v%d%v", par_op, immInt, par_cl)
			}
		} else if op.OpCode.Immediates > 1 {
			immediate := op.Immediates[0].Immediate
			imm := int64(binary.BigEndian.Uint16(immediate))
			if imm > 32767 {
				imm = ((65535 - imm) + 1) * -1
			}
			result = result + par_op + strconv.Itoa(int(imm)) + par_cl
		}

		result += " "
	}

	return result

}

func Evm2Mnem(opcalls []OpCall) string {
	return Evm2Asm(opcalls, "", []string{"(", ")"})
}

func Evm2PyAsm(opcalls []OpCall) string {
	return Evm2Asm(opcalls, "Op.", []string{"[", "]"})
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

			immBytes := make([]byte, len(values)*2+1)
			immBytes[0] = byte(max_size + 1)

			for i, ro := range values {
				relativeOffsetBytes := make([]byte, 2)
				relativeOffset, err := strconv.ParseInt(ro, 10, 64)

				if err != nil {
					return opCall, err
					/*
						fmt.Println("??? err: ", err)
						return opCall, err
						if relativeOffset < 0 {
							relativeOffset = ((65535 - relativeOffset) + 1) * -1
						}
					*/
				}
				binary.BigEndian.PutUint16(relativeOffsetBytes, uint16(relativeOffset))
				immBytes[(i*2)+1] = relativeOffsetBytes[0]
				immBytes[(i*2)+2] = relativeOffsetBytes[1]
			}
			opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immBytes})
			return opCall, nil
		}

		// For >push2 a hexadecimal must be received as parameter
		immBytes := []byte{}
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
			immBytes, err := hex.DecodeString(immediate)
			if err != nil {
				return opCall, err
			}
			opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immBytes})
		} else {
			imm, err := strconv.ParseInt(immediate, 10, 64)
			if err != nil {
				fmt.Println("??? err: ", err)
				immBytes, err = hex.DecodeString(immediate)
				if err != nil {
					return opCall, err
				}
				fmt.Println("??? opcode: ", op.Name)
				fmt.Println("??? immBytes: ", immBytes)
				opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immBytes})
			} else {
				if imm < 0 {
					if op.Name != "RJUMP" && op.Name != "RJUMPI" {
						return opCall, errors.New("Negative immediate only possible for RJUMP and RJUMPI")
					}

					immBytes := make([]byte, 2)
					binary.BigEndian.PutUint16(immBytes, uint16(imm))
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immBytes})
				} else {
					immInt := byte(imm)
					immBytes := make([]byte, op.Immediates)

					if op.Immediates == 2 {
						immInt = byte(imm >> 8)
						immBytes[1] = immInt
					}
					immBytes[0] = immInt
					opCall.Immediates = append(opCall.Immediates, Immediate{Type: Value, Immediate: immBytes})
				}
			}
		}
	}
	return opCall, nil
}

func Mnem2Evm(mn string) (string, error) {
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
				return "", err
			}
			if value < 1 {
				return "", errors.New("Invalid multiplier")
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
			if elements[1][len(elements[1])-1] != ')' {
				return "", errors.New("Invalid immediate: " + token)
			}
		} else {
			opcode = token
		}
		opCall, err := opcode2evm(opcode, immediate)
		if err != nil {
			return "", err
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
				if p, ok := labels[im.Label]; ok {
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
				result += hex.EncodeToString(im.Immediate)
			}
		}
	}
	return result, nil
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

func ModifyEOFObject(eofObject EOFObject, modifier EOFObjectModifier) []byte {
	newcode := make([]byte, 0)
	if modifier.Magic {
		newcode = append(newcode, 0xef)
		newcode = append(newcode, 0x00)
	}
	if modifier.Version {
		newcode = append(newcode, byte(eofObject.Version))
	}
	if modifier.TypeHeader {
		newcode = append(newcode, 0x01)
	}
	if modifier.CodeHeader {
		newcode = append(newcode, 0x02)
		newcode = append(newcode, byte(len(eofObject.CodeSections)))

		for _, cs := range eofObject.CodeSections {
			newcode = append(newcode, byte(len(cs)))
		}
	}
	if modifier.DataHeader {
		newcode = append(newcode, 0x03)
		newcode = append(newcode, byte(len(eofObject.Data)))
	}
	if modifier.Terminator {
		newcode = append(newcode, 0x00)
	}

	for i, _ := range eofObject.Types {
		if s, ok := modifier.TypeSection[i]; ok {
			// If it is blank, means has to be removed
			if s != "" {
				types := []byte(s)
				newcode = append(newcode, types...)
			}
			//} else { // TODO
			//newcode = append(newcode, t...)
		}
	}

	for i, cs := range eofObject.CodeSections {
		if remain, ok := modifier.CodeSection[i]; ok {
			if remain {
				newcode = append(newcode, cs...)
			}
		} else {
			newcode = append(newcode, cs...)
		}
	}

	if modifier.DataSection {
		newcode = append(newcode, eofObject.Data...)
	}
	return newcode
}
