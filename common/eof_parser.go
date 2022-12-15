package common

import (
	"fmt"
	"strconv"
	"strings"
)

type EOFObject struct {
	Version      int64
	Types        [][]int64
	CodeSections []string
	Data         string
}

func (eof *EOFObject) Code() string {
	eof_code := "ef"

	versionHex := fmt.Sprintf("%04x", eof.Version)

	typesHeader := ""
	if len(eof.Types) > 0 {
		typesLengthHex := fmt.Sprintf("%04x", len(eof.Types)*2)
		typesHeader = "03" + typesLengthHex
	}

	codeHeaders := ""
	for _, c := range eof.CodeSections {
		codeLengthHex := fmt.Sprintf("%04x", len(c)/2)
		codeHeaders = codeHeaders + "01" + codeLengthHex
	}

	dataHeader := ""
	if len(eof.Data) > 0 {
		dataLengthHex := fmt.Sprintf("%04x", len(eof.Data)/2)
		dataHeader = dataHeader + "02" + dataLengthHex
	}

	terminator := "00"

	typeContents := ""
	fmt.Println(">types: ", eof.Types)
	for _, t := range eof.Types {
		fmt.Println("t: ", t)
		inputsHex := fmt.Sprintf("%02x", t[0])
		outputsHex := fmt.Sprintf("%02x", t[1])
		typeContents = typeContents + inputsHex + outputsHex
	}
	fmt.Println(">tc: ", typeContents)

	codeContents := ""
	for _, c := range eof.CodeSections {
		codeContents = codeContents + c
	}

	eof_code = eof_code + versionHex + typesHeader + codeHeaders + dataHeader + terminator + typeContents + codeContents + eof.Data
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

func ParseEOF(eof_code string) EOFObject {
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
