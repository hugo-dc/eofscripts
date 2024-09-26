package common

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
)

type EOFObject struct {
	Version           int64
	Types             [][]int64
	CodeSections      [][]byte
	Data              []byte
	ContainerSections [][]byte
	InitCode          bool
	ExplicitRuntime   bool
}

const (
	CTypeId       = 1
	CCodeId       = 2
	CContainerId  = 3
	CDataId       = 4
	CTerminatorId = 0
)

func NewEOFObject() EOFObject {
	return EOFObject{
		Version:           int64(1),
		CodeSections:      make([][]byte, 0),
		Types:             make([][]int64, 0),
		ContainerSections: make([][]byte, 0),
		Data:              make([]byte, 0),
		InitCode:          false,
		ExplicitRuntime:   false,
	}
}

func (eof *EOFObject) EqualTo(other EOFObject) bool {
	if eof.Version != other.Version {
		return false
	}
	if len(eof.Types) != len(other.Types) {
		return false
	}
	for i, t := range eof.Types {
		if t[0] != other.Types[i][0] || t[1] != other.Types[i][1] || t[2] != other.Types[i][2] {
			return false
		}
	}
	if len(eof.CodeSections) != len(other.CodeSections) {
		return false
	}
	for i, c := range eof.CodeSections {
		if string(c) != string(other.CodeSections[i]) {
			return false
		}
	}
	if len(eof.ContainerSections) != len(other.ContainerSections) {
		return false
	}
	for i, c := range eof.ContainerSections {
		if string(c) != string(other.ContainerSections[i]) {
			return false
		}
	}
	if string(eof.Data) != string(other.Data) {
		return false
	}
	if eof.InitCode != other.InitCode {
		return false
	}
	if eof.ExplicitRuntime != other.ExplicitRuntime {
		return false
	}
	return true
}

func (eof *EOFObject) Code() string {
	eof_code := "ef00"

	typeId := fmt.Sprintf("%02x", CTypeId)
	codeId := fmt.Sprintf("%02x", CCodeId)
	containerId := fmt.Sprintf("%02x", CContainerId)
	dataId := fmt.Sprintf("%02x", CDataId)

	versionHex := fmt.Sprintf("%02x", eof.Version)
	typesHeader := ""
	if len(eof.Types) > 0 {
		typesLengthHex := ""
		typesLengthHex = fmt.Sprintf("%04x", len(eof.Types)*4)
		typesHeader = typeId + typesLengthHex
	}

	codeHeaders := ""
	codeLengths := ""
	numCodeSections := 0
	for _, c := range eof.CodeSections {
		codeLengthHex := fmt.Sprintf("%04x", len(c))
		codeLengths = codeLengths + codeLengthHex
		numCodeSections += 1
	}
	numCodeSectionsHex := fmt.Sprintf("%04x", numCodeSections)
	codeHeaders = codeId + numCodeSectionsHex + codeLengths

	containerHeader := ""
	containersLength := ""
	numContainers := 0
	containerContents := ""
	for _, c := range eof.ContainerSections {
		containersLengthHex := fmt.Sprintf("%04x", len(c))
		containersLength = containersLength + containersLengthHex
		numContainers += 1
		containerContents = containerContents + hex.EncodeToString(c)
	}
	if numContainers > 0 {
		numContainersHex := fmt.Sprintf("%04x", numContainers)
		containerHeader = containerId + numContainersHex + containersLength
	}

	dataHeader := ""
	dataLengthHex := fmt.Sprintf("%04x", len(eof.Data))
	dataHeader = dataHeader + dataId + dataLengthHex

	terminator := "00"

	typeContents := ""
	for i, t := range eof.Types {
		inputsHex := fmt.Sprintf("%02x", t[0])
		outputsHex := fmt.Sprintf("%02x", t[1])

		maxStackHeight, isNRF := calculateMaxStackAndNRF(i, eof.CodeSections[i], eof.Types)

		if isNRF { // Mark function as non-returning function
			outputsHex = fmt.Sprintf("%02x", 0x80)
		}

		maxStackHeightHex := fmt.Sprintf("%04x", maxStackHeight)
		typeContents = typeContents + inputsHex + outputsHex + maxStackHeightHex
	}

	codeContents := ""
	for _, c := range eof.CodeSections {
		codeContents = codeContents + hex.EncodeToString(c)
	}

	eof_code = eof_code + versionHex + typesHeader + codeHeaders + containerHeader + dataHeader + terminator + typeContents + codeContents + containerContents + hex.EncodeToString(eof.Data)
	return eof_code
}

func (eof *EOFObject) AddData(dt []byte) {
	eof.Data = dt
}

func (eof *EOFObject) AddCode(code []byte) {
	eof.AddCodeWithType(code, []int64{0, 0x80, 0})
}

func (eof *EOFObject) AddCodeWithType(code []byte, codeType []int64) {
	// Add default type for existing code
	if len(eof.CodeSections) > 0 && len(eof.Types) == 0 {
		eof.Types = append(eof.Types, []int64{0, 0})
	}

	// Add default type for new code section
	eof.Types = append(eof.Types, codeType)

	// Add Code
	eof.CodeSections = append(eof.CodeSections, code)
}

func (eof *EOFObject) AddContainer(container []byte) {
	eof.ContainerSections = append(eof.ContainerSections, container)
}

func (eof *EOFObject) AddDefaultType() bool {
	if len(eof.Types) == 0 {
		eof.Types = append(eof.Types, []int64{0, 0})
		return true
	} else {
		return false
	}
}

func (eof *EOFObject) SetInitcode(initcode bool) {
	eof.InitCode = initcode
}

func (eof *EOFObject) SetExplicitRuntime(explicitRuntime bool) {
	eof.ExplicitRuntime = explicitRuntime
}

func RawBytecodeToSimpleFormat(code []byte) (string, error) {
	ops, err := BytecodeToOpCalls(code)
	code_desc := ""

	if err != nil {
		return code_desc, err
	}

	for _, opc := range ops {
		bc, err := opc.ToBytecode()
		if err != nil {
			panic(err)
		}

		asm := Evm2Mnem([]OpCall{opc})
		code_desc += fmt.Sprintf("%6x # [%v] %v\n", bc, opc.Position, asm)
	}
	return code_desc, nil
}

func RawBytecodeToPythonFormat(code []byte) (string, error) {
	ops, err := BytecodeToOpCalls(code)
	if err != nil {
		return "", err
	}

	prev_op := OpCall{}
	op_count := 1
	descriptions := []string{}
	code_desc := ""
	for _, opc := range ops {
		if opc.EqualTo(prev_op) {
			op_count++
			continue
		}
		asm := Evm2PyAsm([]OpCall{opc})
		if op_count > 1 {
			code_desc += fmt.Sprintf("* %v ", op_count)
			descriptions = append(descriptions, code_desc)
			code_desc = fmt.Sprintf("%v", asm)
			op_count = 1
		} else {
			if code_desc != "" {
				descriptions = append(descriptions, code_desc)
			}
			code_desc = fmt.Sprintf("%v", asm)
		}
		prev_op = opc
	}
	if op_count > 1 {
		code_desc += fmt.Sprintf("* %v ", op_count)
	}
	code_desc = strings.TrimRight(code_desc, " ")
	descriptions = append(descriptions, code_desc)
	return strings.Join(descriptions, "+ "), nil
}

func (eof *EOFObject) Describe() string {
	code_section_headers := ""
	for i, v := range eof.CodeSections {
		code_section_headers += fmt.Sprintf("  %04x ", len(v)) + fmt.Sprintf("# Code section %v, %v bytes\n", i, len(v))
	}

	container_section_headers := ""
	if len(eof.ContainerSections) > 0 {
		container_section_headers = fmt.Sprintf("%02x%04x # Total container sections (%v)\n", CContainerId, len(eof.ContainerSections), len(eof.ContainerSections))
		for i, v := range eof.ContainerSections {
			container_section_headers += fmt.Sprintf("  %04x # Container section %v, %v bytes\n", len(v), i, len(v))
		}
	}
	if container_section_headers == "" {
		container_section_headers = "       # No container sections"
	}

	types_headers := ""
	for i, v := range eof.Types {
		types_headers += fmt.Sprintf("       # Code %v types\n", i) +
			fmt.Sprintf("    %02x", v[0]) + fmt.Sprintf(" # %v inputs\n", v[0])
		if v[1] == 0x80 {
			types_headers += fmt.Sprintf("    %02x", v[1]) + " # 0 outputs (Non-returning function)\n"
		} else {
			types_headers += fmt.Sprintf("    %02x", v[1]) + fmt.Sprintf(" # %v outputs\n", v[1])
		}
		types_headers += fmt.Sprintf("  %04x", v[2]) + fmt.Sprintf(" # max_stack: %v\n", v[2])
	}

	code_sections := ""
	for i, v := range eof.CodeSections {
		code_sections += fmt.Sprintf("       # Code section %v\n", i)
		code_description, err := RawBytecodeToSimpleFormat(v)
		if err != nil {
			code_sections += hex.EncodeToString(v)
		} else {
			code_sections += code_description
		}
	}

	container_sections := ""
	if len(eof.ContainerSections) > 0 {
		container_sections += "       # Container sections\n"
		for i, v := range eof.ContainerSections {
			container_sections += fmt.Sprintf("       #   Container section %v\n", i)
			code_description, err := RawBytecodeToSimpleFormat(v)
			if err != nil {
				container_sections += hex.EncodeToString(v) + "\n"
			} else {
				container_sections += code_description
			}
		}
	}

	comment := ""
	if len(eof.Data) == 0 {
		comment = "(empty)"
	}

	data_section := fmt.Sprintf("       # Data section %s\n", comment)
	data_section += hex.EncodeToString(eof.Data)

	eof_desc := ""
	if len(eof.ContainerSections) == 0 {
		eof_desc = `EF00%02x # Magic and Version (%v)
%02x%04x # Types length (%v)
%02x%04x # Total code sections (%v)
%s%02x%04x # Data section length (%v)
    %02x # Terminator (end of header)
%s%s%s
`
		eof_desc = fmt.Sprint(fmt.Sprintf(eof_desc, eof.Version, eof.Version, CTypeId, len(eof.Types)*4, len(eof.Types)*4, CCodeId, len(eof.CodeSections), len(eof.CodeSections), code_section_headers, CDataId, len(eof.Data)/2, len(eof.Data)/2, CTerminatorId, types_headers, code_sections, data_section))
	} else {
		eof_desc = `EF00%02x # Magic and Version (%v)	
%02x%04x # Types length (%v)
%02x%04x # Total code sections (%v)
%s%s%02x%04x # Data section length (%v)
%s%s%s%s
`

		eof_desc = fmt.Sprint(fmt.Sprintf(eof_desc, eof.Version, eof.Version, CTypeId, len(eof.Types)*4, len(eof.Types)*4, CCodeId, len(eof.CodeSections), len(eof.CodeSections), code_section_headers, container_section_headers, CDataId, len(eof.Data)/2, len(eof.Data)/2, types_headers, code_sections, container_sections, data_section))
	}
	return eof_desc
}

func (eof *EOFObject) DescribeAsPython(depth uint16, index uint16) string {
	pyeof_code := ""
	str_depth := ""
	if depth != 0 {
		str_depth = fmt.Sprintf("_D%vI%v", depth, index)
	}

	if len(eof.ContainerSections) == 0 {
		pyeof_code = `Container(
  name = 'EOFV1_0000%s',
  sections = [
    %s%s%s  ],%s
)`
	} else {
		pyeof_code = `Container(
  name = 'EOFV1_0000%s',
  sections = [
    %s%s%s  
  ],%s
)`
	}
	code_sections := ""
	for i, v := range eof.CodeSections {
		// Get Max Stack Height
		max_stack_height := eof.Types[i][2]
		code := ""
		code_description, err := RawBytecodeToPythonFormat(v)
		if err != nil {
			code = hex.EncodeToString(v) + "\n"
		} else {
			code = code_description
		}
		extra_args := ""
		if eof.Types[i][0] != 0 {
			extra_args = fmt.Sprintf(" code_inputs=%v,", eof.Types[i][0])
		}
		if eof.Types[i][1] != 0x80 {
			extra_args += fmt.Sprintf(" code_outputs=%v,", eof.Types[i][1])
		}

		code_sections += fmt.Sprintf("  Section.Code(code=%s,%s max_stack_height=%v),\n    ", code, extra_args, max_stack_height)
	}

	container_sections := ""
	if len(eof.ContainerSections) > 0 {
		for i, v := range eof.ContainerSections {
			raw_bytecode := ""
			for i := 0; i < len(v); i += 2 {
				separator := ", "
				if i == len(v)-2 {
					separator = ""
				}
				raw_bytecode += fmt.Sprintf("0x%x%s ", v[i], separator)
			}

			formatted_subcontainer, error := ParseEOF(v)
			if error != nil {
				subcontainerId := fmt.Sprintf("_D%vI%v", depth+1, i)
				container_sections += fmt.Sprintf(`  Section.Container(
          container=Container(
              name="EOFV1_0000%s",
              raw_bytes=bytes(
                  [ %s ])
          )
      ),
    `, subcontainerId, raw_bytecode)
			} else {
				subcontainer := formatted_subcontainer.DescribeAsPython(depth+1, uint16(i))
				subcontainer = strings.Replace(subcontainer, "\n", "\n        ", -1)
				container_sections += fmt.Sprintf(`  Section.Container(container=%s
      ),`, subcontainer)
			}
		}
	}

	data_section := ""
	if len(eof.Data) > 0 {
		data_section += fmt.Sprintf("  Section.Data(data=\"%s\")\n", eof.Data)
	}

	container_kind := ""
	if eof.ExplicitRuntime {
		container_kind = "\n  kind=ContainerKind.RUNTIME"
	}
	if eof.InitCode {
		container_kind = "\n  kind=ContainerKind.INITCODE"
	}

	return fmt.Sprintf(pyeof_code, str_depth, code_sections, container_sections, data_section, container_kind)
}

func calculateMaxStackAndNRF(funcId int, code []byte, types [][]int64) (int64, bool) {
	stackHeights := map[int64]int64{}
	startStackHeight := types[funcId][0]
	maxStackHeight := startStackHeight
	isNRF := true
	worklist := [][]int64{{0, startStackHeight}}

	opCodes := GetOpcodesByNumber()

	for {
		ix := len(worklist) - 1
		res := worklist[ix]
		pos := res[0]
		stackHeight := res[1]
		worklist = worklist[:ix]

	outer:
		for int(pos)+1 < len(code) {
			if pos < 0 {
				fmt.Println("Position out of bounds: ", pos)
				break
			}
			op := uint16(binary.BigEndian.Uint16([]byte{code[pos], code[pos+1]}))
			opCode := opCodes[int(op)]
			if exp, ok := stackHeights[pos]; ok {
				if stackHeight != exp {
					fmt.Println("stackHeight:", stackHeight, "exp:", exp, "at pos", pos)
					fmt.Println("Error: stack height mismatch for different paths")
					break
				} else {
					break
				}
			} else {
				stackHeights[pos] = stackHeight
			}

			stackHeightRequired := int64(opCode.StackInput)
			stackHeightChange := int64(opCode.StackOutput - opCode.StackInput)

			if stackHeightRequired > stackHeight {
				fmt.Println("stackHeightRequired:", stackHeightRequired, "stackHeight:", stackHeight)
				fmt.Println("stack underflow")
			}

			if (1024 + stackHeightChange) < stackHeight {
				fmt.Println("stack overflow")
				break
			}

			stackHeight += stackHeightChange
			switch {
			case opCode.Name == "CALLF":
				if int(pos+3) > len(code) {
					fmt.Println("truncated CALLF")
					break
				}
				calledFuncId := int64(binary.BigEndian.Uint16([]byte{code[pos+2], code[pos+3]}))
				if int(calledFuncId) >= len(types) {
					fmt.Println("invalid function id")
					break
				}
				stackHeightRequired += int64(types[calledFuncId][0])
				stackHeightChange += int64(types[calledFuncId][1] - types[calledFuncId][0])

				if stackHeightRequired > stackHeight {
					fmt.Println("stack underflow")
					break
				}

				if (1024 + stackHeightChange) < stackHeight {
					fmt.Println("stack overflow")
					break
				}
				stackHeight += stackHeightChange
				pos += 3
			case opCode.Name == "RETF":
				isNRF = false
				if int64(types[funcId][1]) != stackHeight {
					fmt.Printf("Wrong number of outputs (want:%d, got: %d)\n", types[funcId][1], stackHeight)
				}
				break outer
			case opCode.Name == "RJUMP":
				if int(pos)+3 > len(code) {
					fmt.Println("Error: Truncted RJUMP")
					break
				}
				offset := int64(binary.BigEndian.Uint16([]byte{code[pos+1], code[pos+3]}))
				pos += (int64(opCode.Immediates) + 1 + int64(offset))
			case opCode.Name == "RJUMPI":
				if int(pos)+3 > len(code) {
					fmt.Println("Error: Truncted RJUMPI")
					break
				}

				offset := int64(binary.BigEndian.Uint16([]byte{code[pos+1], code[pos+3]}))

				if offset > 32767 {
					offset = ((65535 - offset) + 1) * -1
				}

				worklist = append(worklist, []int64{pos + 3 + offset, stackHeight})
				pos += int64(opCode.Immediates) + 1
			case opCode.Name == "RJUMPV":
				if int(pos)+3 > len(code) {
					fmt.Println("Error: truncated RJUMPV")
					break
				}
				count := int64(binary.BigEndian.Uint16([]byte{code[pos+1], code[pos+2]}))
				count = count + 1
				fmt.Println("\tcount:", count)

				pos += 2
				if int(pos)+int(count) > len(code) {
					fmt.Println("Error: truncated RJUMPV")
					break
				}
				for i := 0; i < int(count); i++ {
					fmt.Println("codelen:", len(code))
					fmt.Println("target: ", int(pos)+2*i+2)
					if len(code) <= int(pos)+2*i+2 {
						fmt.Println("Error: truncated RJUMPV")
						break
					}

					offset := int64(binary.BigEndian.Uint16([]byte{code[pos+2*int64(i)], code[pos+2*int64(i)+2]}))
					if offset > 32767 {
						offset = ((65535 - offset) + 1) * -1
					}

					fmt.Println("\toffset:", offset)
					fmt.Println("\twE:", pos+count+offset)
					worklist = append(worklist, []int64{pos + count + offset, stackHeight})
				}
				pos += count
			default:
				if opCode.IsTerminating {
					break outer
				} else {
					pos += int64(opCode.Immediates) + 1
				}
			}
			maxStackHeight = int64(math.Max(float64(maxStackHeight), float64(stackHeight)))
		}

		if maxStackHeight > 1024 {
			fmt.Println("Error: max stack above limit")
		}

		if len(worklist) == 0 {
			break
		}
	}
	return maxStackHeight, isNRF
}

func consumeMagicAndVersion(bytecode []byte) ([]byte, int64, error) {
	if len(bytecode) < 3 {
		return bytecode, 0, errors.New("Invalid EOF code")
	}

	if bytecode[0] != 0xef {
		return bytecode, 0, errors.New("Invalid EOF code")
	}

	versionHex := bytecode[1:3]
	version := int64(binary.BigEndian.Uint16(versionHex))
	return bytecode[3:], version, nil
}

func consumeTypesHeader(bytecode []byte) ([]byte, int64, error) {
	if len(bytecode) < 3 {
		return bytecode, 0, errors.New("Invalid types header")
	}

	typesLengthBytes := bytecode[1:3]
	typesLength := int64(binary.BigEndian.Uint16(typesLengthBytes))
	return bytecode[3:], typesLength, nil
}

func consumeCodeHeader(bytecode []byte, container bool) ([]byte, []int64, error) {
	codeHeaders := []int64{}
	if len(bytecode) < 5 {
		return bytecode, []int64{}, errors.New("Invalid code section (1)")
	}
	if bytecode[0] != 0x02 && !container {
		return bytecode, []int64{}, errors.New("Invalid code section (2)")
	}
	if bytecode[0] != 0x03 && container {
		return bytecode, []int64{}, errors.New("Invalid container section")
	}

	codeSectionLengthBytes := bytecode[1:3]
	codeSectionLength := int64(binary.BigEndian.Uint16(codeSectionLengthBytes))

	for i := 3; i < int(codeSectionLength)+3; i++ {
		codeLengthBytes := bytecode[i : i+2]
		codeLength := int64(binary.BigEndian.Uint16(codeLengthBytes))
		codeHeaders = append(codeHeaders, codeLength)
		i += 1
	}

	return bytecode[int(codeSectionLength)+4:], codeHeaders, nil
}

func consumeDataSection(bytecode []byte) ([]byte, []byte, error) {
	if len(bytecode) < 3 {
		return bytecode, []byte{}, errors.New("Invalid data section")
	}

	if bytecode[0] != 0x04 {
		return bytecode, []byte{}, errors.New("Invalid data section")
	}

	dataLengthBytes := bytecode[1:3]
	dataLength := int64(binary.BigEndian.Uint16(dataLengthBytes))

	if dataLength == 0 {
		return bytecode[3:], []byte{}, nil
	}
	return bytecode[3:], bytecode[3 : 3+dataLength], nil
}

func consumeTerminator(bytecode []byte) ([]byte, int64, error) {
	if len(bytecode) < 1 {
		return bytecode, 0, errors.New("Invalid terminator")
	}

	if bytecode[0] != 0x00 {
		return bytecode, 0, errors.New("Invalid terminator")
	}

	return bytecode[1:], 0, nil
}

func consumeTypesContent(bytecode []byte, typesLength int64) ([]byte, [][]int64, error) {
	types := [][]int64{}
	i := 0
	for {
		if i >= int(typesLength) {
			break
		}

		inputs := int64(bytecode[i])
		outputs := int64(bytecode[i+1])

		maxStackBytes := bytecode[i+2 : i+4]
		maxStack := int64(binary.BigEndian.Uint16(maxStackBytes))
		types = append(types, []int64{inputs, outputs, maxStack})
		i += 4
	}
	return bytecode[i:], types, nil
}

func consumeCodeContent(bytecode []byte, codeHeaders []int64) ([]byte, [][]byte, error) {
	codeSections := [][]byte{}
	i := 0
	for {
		if i >= len(codeHeaders) {
			break
		}
		code := bytecode[i : i+int(codeHeaders[i])]
		codeSections = append(codeSections, code)
		i += int(codeHeaders[i])
	}
	return bytecode[i:], codeSections, nil
}

func ParseEOF(eof_code []byte) (EOFObject, error) {
	typesLength := int64(0)
	codeHeaders := []int64{}
	containerHeaders := []int64{}

	result := NewEOFObject()
	err := errors.New("")

	eof_code, result.Version, err = consumeMagicAndVersion(eof_code)
	if err != nil {
		return result, err
	}

	eof_code, typesLength, err = consumeTypesHeader(eof_code)
	if err != nil {
		return result, err
	}

	eof_code, codeHeaders, err = consumeCodeHeader(eof_code, false)
	if err != nil {
		return result, err
	}
	eof_code, containerHeaders, err = consumeCodeHeader(eof_code, true)

	eof_code, result.Data, err = consumeDataSection(eof_code)
	if err != nil {
		return result, err
	}

	eof_code, _, err = consumeTerminator(eof_code)
	if err != nil {
		return result, err
	}

	eof_code, result.Types, err = consumeTypesContent(eof_code, typesLength)
	if err != nil {
		return result, err
	}

	eof_code, result.CodeSections, err = consumeCodeContent(eof_code, codeHeaders)
	if err != nil {
		return result, err
	}

	eof_code, result.ContainerSections, err = consumeCodeContent(eof_code, containerHeaders)
	if err != nil {
		return result, err
	}

	if len(eof_code) > 0 {
		return result, errors.New("EOF code not fully consumed")
	}

	return result, nil
}
