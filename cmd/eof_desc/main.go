package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/ethscripts/common"
)

func describeCode(code string) {
	ops, err := common.DescribeBytecode(code)

	if err != nil {
		panic(fmt.Sprintf("%v: code %s", err, code))
	}

	for _, opc := range ops {
		bc, err := opc.ToBytecode()

		if err != nil {
			panic(err)
		}

		asm := common.Evm2Mnem(bc)
		fmt.Println(fmt.Sprintf("% 6v # [%v] %v", bc, opc.Position, asm))
	}
}

func main() {
	eof_code := os.Args[1]
	if eof_code[:2] == "0x" {
		eof_code = eof_code[2:]
	}
	eof_code = strings.Replace(eof_code, " ", "", -1)
	eofObject, err := common.ParseEOF(eof_code)

	if err != nil {
		panic(err)
	}

	// Print the EOF object
	fmt.Println("EF00"+fmt.Sprintf("%02x", eofObject.Version), "# Magic and Version (", eofObject.Version, ")")
	fmt.Println(fmt.Sprintf("%02x%04x", common.CTypeId, len(eofObject.Types)*4), "# Types length (", len(eofObject.Types)*4, ")")
	fmt.Println(fmt.Sprintf("%02x%04x", common.CCodeId, len(eofObject.CodeSections)), "# Total code sections (", len(eofObject.CodeSections), ")")

	for i, v := range eofObject.CodeSections {
		fmt.Println(fmt.Sprintf("  %04x", len(v)/2), "# Code section ", i, ",", len(v)/2, " bytes")
	}

	if len(eofObject.ContainerSections) > 0 {
		fmt.Println(fmt.Sprintf("%02x%04x", common.CContainerId, len(eofObject.ContainerSections)), "# Total container sections (", len(eofObject.ContainerSections), ")")
		for i, v := range eofObject.ContainerSections {
			fmt.Println(fmt.Sprintf("  %04x", len(v)/2), "# Container section ", i, ",", len(v)/2, " bytes")
		}
	}
	fmt.Println(fmt.Sprintf("%02x%04x", common.CDataId, len(eofObject.Data)/2), "# Data section length (", len(eofObject.Data)/2, ")")
	fmt.Println(fmt.Sprintf("    %02x # Terminator (end of header)", common.CTerminatorId))

	for i, v := range eofObject.Types {
		fmt.Println("       # Code", i, "types")
		fmt.Println(fmt.Sprintf("    %02x", v[0]), "#", v[0], "inputs")
		if v[1] == 0x80 {
			fmt.Println(fmt.Sprintf("    %02x", v[1]), "# 0 outputs", "(Non-returning function)")
		} else {
			fmt.Println(fmt.Sprintf("    %02x", v[1]), "#", v[1], "outputs")
		}
		fmt.Println(fmt.Sprintf("  %04x", v[2]), "#", "max stack:", v[2])
	}

	for i, v := range eofObject.CodeSections {
		fmt.Println("       # Code section", i)
		describeCode(v)
	}

	if len(eofObject.ContainerSections) > 0 {
		fmt.Println("       # Container sections")
		for i, v := range eofObject.ContainerSections {
			fmt.Println("       # Container section", i)
			fmt.Println(v)
			//describeCode(v)
		}
	}

	comment := ""
	if len(eofObject.Data) == 0 {
		comment = "(empty)"
	}

	fmt.Println("       # Data section", comment)
	fmt.Println(eofObject.Data)
	/*
		for i, v := range eofObject.ContainerSections {
			fmt.Println("       # Container section", i)
			fmt.Println(v)
		}
	*/
}
