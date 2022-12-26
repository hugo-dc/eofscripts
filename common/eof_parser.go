package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type EOFObject struct {
	Version      int64
	Types        [][]int64
	CodeSections []string
	Data         string
}

const (
	cOldCodeId = "01"
	cOldDataId = "02"
	cOldTypeId = "03"
	cTypeId    = "01"
	cCodeId    = "02"
	cDataId    = "03"
)

func NewEOFObject() EOFObject {
	return EOFObject{
		Version:      int64(1),
		CodeSections: make([]string, 0),
		Types:        make([][]int64, 0),
		Data:         "",
	}
}

func (eof *EOFObject) CodeNew(withTypes bool) string {
	return eof.Code(false, withTypes)
}

func (eof *EOFObject) Code(old bool, withTypes bool) string {
	eof_code := "ef00"

	typeId := cTypeId
	codeId := cCodeId
	dataId := cDataId

	if old {
		typeId = cOldTypeId
		codeId = cOldCodeId
		dataId = cOldDataId
	}

	versionHex := fmt.Sprintf("%02x", eof.Version)
	typesHeader := ""
	if len(eof.Types) > 0 {
		typesLengthHex := fmt.Sprintf("%04x", len(eof.Types)*2)
		typesHeader = typeId + typesLengthHex
	}

	codeHeaders := ""
	oldCodeHeaders := ""
	codeLengths := ""
	numCodeSections := 0
	for _, c := range eof.CodeSections {
		codeLengthHex := fmt.Sprintf("%04x", len(c)/2)
		codeLengths = codeLengths + codeLengthHex
		oldCodeHeaders = oldCodeHeaders + codeId + codeLengthHex
		numCodeSections += 1
	}
	numCodeSectionsHex := fmt.Sprintf("%04x", numCodeSections)
	codeHeaders = codeId + numCodeSectionsHex + codeLengths

	dataHeader := ""
	dataLengthHex := fmt.Sprintf("%04x", len(eof.Data)/2)
	dataHeader = dataHeader + dataId + dataLengthHex

	terminator := "00"

	typeContents := ""
	for i, t := range eof.Types {
		inputsHex := fmt.Sprintf("%02x", t[0])
		outputsHex := fmt.Sprintf("%02x", t[1])

		maxStackHeight := calculateMaxStack(i, eof.CodeSections[i], eof.Types)
		maxStackHeightHex := fmt.Sprintf("%04x", maxStackHeight)
		typeContents = typeContents + inputsHex + outputsHex + maxStackHeightHex
	}

	codeContents := ""
	for _, c := range eof.CodeSections {
		codeContents = codeContents + c
	}

	if withTypes == false && len(eof.Types) == 1 && old == true {
		/*
			fmt.Println("magic:", eof_code)
			fmt.Println("version:", versionHex)
			fmt.Println("oldCodeHeaders:", oldCodeHeaders)
			fmt.Println("dataHeader:", dataHeader)
			fmt.Println("terminator:", terminator)
			fmt.Println("codeContents:", codeContents)
			fmt.Println("eofData:", eof.Data)
		*/
		if len(eof.Data) > 0 {
			eof_code = eof_code + versionHex + oldCodeHeaders + dataHeader + terminator + codeContents + eof.Data
		} else {
			eof_code = eof_code + versionHex + oldCodeHeaders + terminator + codeContents
		}
	} else {
		fmt.Println("HEADER\n--------")
		fmt.Println("magic:", eof_code)
		fmt.Println("version:", versionHex)
		fmt.Println("types:", typesHeader)
		fmt.Println("codeHeaders:", codeHeaders)
		fmt.Println("dataHeader:", dataHeader)
		fmt.Println("terminator:", terminator)
		fmt.Println("BODY\n--------")
		fmt.Println("typesSection:", typeContents)
		fmt.Println("codeSection:", codeContents)
		fmt.Println("dataSection:", eof.Data)
		eof_code = eof_code + versionHex + typesHeader + codeHeaders + dataHeader + terminator + typeContents + codeContents + eof.Data
	}
	return eof_code
}

func (eof *EOFObject) AddData(dt string) {
	eof.Data = dt
}

func (eof *EOFObject) AddCode(code string) {
	eof.AddCodeWithType(code, []int64{0, 0})
}

func (eof *EOFObject) AddCodeWithType(code string, codeType []int64) {
	// Add default type for existing code
	if len(eof.CodeSections) > 0 && len(eof.Types) == 0 {
		eof.Types = append(eof.Types, []int64{0, 0})
	}

	// Add default type for new code section
	eof.Types = append(eof.Types, codeType)

	// Add Code
	eof.CodeSections = append(eof.CodeSections, code)
}

func (eof *EOFObject) AddDefaultType() bool {
	if len(eof.Types) == 0 {
		eof.Types = append(eof.Types, []int64{0, 0})
		return true
	} else {
		return false
	}
}

func calculateMaxStack(funcId int, code string, types [][]int64) int64 {
	stackHeights := map[int64]int64{}
	startStackHeight := types[funcId][0]
	maxStackHeight := startStackHeight
	worklist := [][]int64{{0, startStackHeight}}

	opCodes := GetOpcodesByNumber()

	for {
		ix := len(worklist) - 1
		res := worklist[ix]
		pos := res[0]
		stackHeight := res[1]
		worklist = worklist[:ix]

		for {
			if int(pos) >= (len(code) / 2) {
				fmt.Println("Error: code is invalid")
				return 0
			}

			op, _ := strconv.ParseInt(code[pos*2:pos*2+2], 16, 64)
			opCode := opCodes[int(op)]
			if _, ok := stackHeights[pos]; ok {
				if stackHeight != stackHeights[pos] {
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

			if opCode.Name == "CALLF" {
				calledFuncId, _ := strconv.ParseInt(code[pos+1:pos+3], 16, 64)

				stackHeightRequired += int64(types[calledFuncId][0])
				stackHeightChange += int64(types[calledFuncId][1] - types[calledFuncId][0])
			}

			if stackHeight < stackHeightRequired {
				fmt.Println("Error: stack underflow")
				break
			}

			stackHeight += stackHeightChange
			maxStackHeight = int64(math.Max(float64(maxStackHeight), float64(stackHeight)))

			// jumps
			if opCode.Name == "RJUMP" {
				offset, _ := strconv.ParseInt(code[pos+1:pos+3], 16, 64)
				pos += int64(opCode.Immediates) + 1 + int64(offset)
			} else if opCode.Name == "RJUMPI" {
				offset, _ := strconv.ParseInt(code[pos+1:pos+3], 16, 64)
				worklist = append(worklist, []int64{pos + 3 + offset, stackHeight})
				pos += int64(opCode.Immediates) + 1
			} else if opCode.IsTerminating {
				expectedHeight := 0
				if opCode.Name == "RETF" {
					expectedHeight = int(types[funcId][1])
				}
				if int(stackHeight) != expectedHeight {
					fmt.Println("Warning: Non-empty stack on terminating instruction")
				}
				break
			} else {
				pos += int64(opCode.Immediates) + 1
			}
		}

		if maxStackHeight > 1024 {
			fmt.Println("Error: max stack above limit")
		}

		if len(worklist) == 0 {
			break
		}
	}
	return maxStackHeight
}

func ParseOldEOF(eof_code string) EOFObject {
	version := int64(0)
	versionHex := ""
	eof_code = strings.ToLower(eof_code)

	codeHeaders := []int64{}
	codeSections := []string{}
	typesLength := int64(0)
	types := [][]int64{}
	dataLength := int64(0)
	dataContent := ""

	i := 0
	for {
		if i+2 > len(eof_code) {
			break
		}

		bt := eof_code[i : i+2]
		fmt.Println("bt: ", bt)

		if versionHex == "" && bt != "ef" {
			fmt.Println("Error: Invalid EOF code")
			break
		}

		if versionHex == "" && bt == "ef" {
			versionHex = eof_code[i+2 : i+6]
			version, _ = strconv.ParseInt(versionHex, 16, 64)
			i += 4
		}

		if bt == "03" {
			fmt.Println(">types")
			fmt.Println("code: ", eof_code[i:])
			typesLengthHex := eof_code[i+2 : i+6]
			typesLengthTmp, err := strconv.ParseInt(typesLengthHex, 16, 64)

			if err != nil {
				fmt.Println("Error: Invalid types legnth : ", err)
			}

			typesLength = typesLengthTmp
			i += 4
		}

		if bt == "01" {
			fmt.Println(">code")
			fmt.Println("code: ", eof_code[i:])
			codeLenHex := eof_code[i+2 : i+6]
			codeLen, err := strconv.ParseInt(codeLenHex, 16, 64)

			if err != nil {
				fmt.Println("Error: Invalid code length :", err)
			}

			codeHeaders = append(codeHeaders, codeLen)
			i += 4
		}

		if bt == "02" {
			fmt.Println(">data")
			fmt.Println("code: ", eof_code[i:])
			dataLengthHex := eof_code[i+2 : i+6]
			dataLength, _ = strconv.ParseInt(dataLengthHex, 16, 64)
			i += 4
		}

		if bt == "00" { // Terminator
			fmt.Println("> terminator")
			fmt.Println("code: ", eof_code[i:])
			// Extract Types
			if typesLength > 0 {
				for j := 0; j < int(typesLength); j += 2 {
					inputsHex := eof_code[i+2 : i+4]
					inputs, _ := strconv.ParseInt(inputsHex, 16, 64)
					outputsHex := eof_code[i+4 : i+6]
					outputs, _ := strconv.ParseInt(outputsHex, 16, 64)

					types = append(types, []int64{inputs, outputs})
					i += 4
				}
			}
			fmt.Println("types: ", types)

			i += 2
			fmt.Println("code: ", eof_code[i:])
			fmt.Println("codeHeaders: ", codeHeaders)
			// Extract Code
			for _, cH := range codeHeaders {
				code := eof_code[i : i+int(cH)*2]
				fmt.Println("codeX: ", code)
				codeSections = append(codeSections, code)
				i += int(cH) * 2
			}

			// Extract Data
			dataContent = eof_code[i : i+int(dataLength)*2]
			i += int(dataLength) * 2
		}
		i += 2
	}
	return EOFObject{
		Version:      version,
		CodeSections: codeSections,
		Types:        types,
		Data:         dataContent,
	}
}
