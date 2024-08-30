package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

func describeCode(code string) string {
	ops, err := common.DescribeBytecode(code)
	code_desc := ""

	if err != nil {
		panic(fmt.Sprintf("%v: code %s", err, code))
	}

	for _, opc := range ops {
		bc, err := opc.ToBytecode()
		if err != nil {
			panic(err)
		}

		asm := common.Evm2Mnem([]common.OpCall{opc})
		code_desc += fmt.Sprintf("%6s # [%v] %v\n", bc, opc.Position, asm)
	}
	return code_desc
}

func main() {
	eof_code := ""
	if len(os.Args) < 2 {
		fmt.Scanln(&eof_code)
	} else {
		stdin_input := ""
		fmt.Scanln(&stdin_input)
		if len(stdin_input) > 0 {
			fmt.Println("Error: No arguments allowed when reading from stdin")
			return
		}
		eof_code = strings.Join(os.Args[1:], " ")
	}
	if eof_code[:2] == "0x" {
		eof_code = eof_code[2:]
	}
	eof_code = strings.Replace(eof_code, " ", "", -1)
	eofObject, err := common.ParseEOF(eof_code)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Print the EOF object
	code_section_headers := ""
	for i, v := range eofObject.CodeSections {
		code_section_headers += fmt.Sprintf("  %04x ", len(v)/2) + fmt.Sprintf("# Code section %v, %v bytes\n", i, len(v)/2)
	}

	container_section_headers := ""
	if len(eofObject.ContainerSections) > 0 {
		container_section_headers += fmt.Sprintf("%02x%04x", common.CContainerId, len(eofObject.ContainerSections)) + fmt.Sprintf("# Total container sections (%i)\n", len(eofObject.ContainerSections))
		for i, v := range eofObject.ContainerSections {
			container_section_headers += fmt.Sprintf("  %04x", len(v)/2) + fmt.Sprintf("# Container section %i, $i bytes\n", i, len(v)/2)
		}
	}
	if container_section_headers == "" {
		container_section_headers = "       # No container sections"
	}

	types_headers := ""
	for i, v := range eofObject.Types {
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
	for i, v := range eofObject.CodeSections {
		code_sections += fmt.Sprintf("       # Code section %v\n", i)
		code_sections += describeCode(v)
	}

	container_sections := ""
	if len(eofObject.ContainerSections) > 0 {
		container_sections += "       # Container sections"
		for i, v := range eofObject.ContainerSections {
			container_sections += fmt.Sprintf("       # Container section", i)
			container_sections += v
			//describeCode(v)
		}
	}

	comment := ""
	if len(eofObject.Data) == 0 {
		comment = "(empty)"
	}

	data_section := fmt.Sprintf("       # Data section %s\n", comment)
	data_section += eofObject.Data

	eof_desc := ""
	if len(eofObject.ContainerSections) == 0 {
		eof_desc = `EF00%02x # Magic and Version (%v)
%02x%04x # Types length (%v)
%02x%04x # Total code sections (%v)
%s%02x%04x # Data section length (%v)
    %02x # Terminator (end of header)
%s%s%s
`
		fmt.Print(fmt.Sprintf(eof_desc, eofObject.Version, eofObject.Version, common.CTypeId, len(eofObject.Types)*4, len(eofObject.Types)*4, common.CCodeId, len(eofObject.CodeSections), len(eofObject.CodeSections), code_section_headers, common.CDataId, len(eofObject.Data)/2, len(eofObject.Data)/2, common.CTerminatorId, types_headers, code_sections, data_section))
	} else {
		eof_desc = `EF00%02x # Magic and Version (%v)	
%02x%04x # Types length (%v)
%02x%04x # Total code sections (%v)
%s%s%02x%04x # Total container sections (%v)
%s%s%s
`

		fmt.Print(fmt.Sprintf(eof_desc, eofObject.Version, eofObject.Version, common.CTypeId, len(eofObject.Types)*4, len(eofObject.Types)*4, common.CCodeId, len(eofObject.CodeSections), code_section_headers, container_section_headers, common.CDataId, len(eofObject.Data)/2, len(eofObject.Data)/2, len(eofObject.ContainerSections), types_headers, code_sections, container_sections, data_section))
	}

	/*
		for i, v := range eofObject.ContainerSections {
			fmt.Println("       # Container section", i)
			fmt.Println(v)
		}
	*/
}
