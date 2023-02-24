package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/ethscripts/common"
)

func describeCode(code string) {
	ops, err := common.DescribeBytecode(code)

	if err != nil {
		panic(err)
	}

	pc := 0
	for _, opc := range ops {
		bc, err := opc.ToBytecode()

		if err != nil {
			panic(err)
		}

		asm := common.Evm2Mnem(bc)

		fmt.Println(fmt.Sprintf("% 6v # %v", bc, asm))

		pc += 1 + opc.OpCode.Immediates
	}
}

func main() {
	eof_code := os.Args[1]
	eofObject, err := common.ParseEOF(eof_code)

	if err != nil {
		panic(err)
	}

	// Print the EOF object
	fmt.Println("EF00"+strconv.Itoa(int(eofObject.Version)), " # Magic and Version (", eofObject.Version, ")")
	fmt.Println(fmt.Sprintf("01%04x", len(eofObject.Types)*4), "# Types length (", len(eofObject.Types)*4, ")")
	fmt.Println(fmt.Sprintf("02%04x", len(eofObject.CodeSections)), "# Total code sections (", len(eofObject.CodeSections), ")")

	for i, v := range eofObject.CodeSections {
		fmt.Println(fmt.Sprintf("  %04x", len(v)/2), "# Code section ", i, ",", len(v)/2, " bytes")
	}

	fmt.Println(fmt.Sprintf("03%04x", len(eofObject.Data)), "# Data section lenth (", len(eofObject.Data), ")")
	fmt.Println("    00", "# Terminator (end of header)")

	for i, v := range eofObject.Types {
		fmt.Println("       # Code", i, "types")
		fmt.Println(fmt.Sprintf("    %02x", v[0]), "#", i, "inputs")
		fmt.Println(fmt.Sprintf("    %02x", v[1]), "#", i, "outputs")
		fmt.Println(fmt.Sprintf("  %04x", v[2]), "#", "max stack:", v[2])
	}

	for i, v := range eofObject.CodeSections {
		fmt.Println("       # Code section", i)
		describeCode(v)
	}

	comment := ""

	if len(eofObject.Data) == 0 {
		comment = "(empty)"
	}

	fmt.Println("       # Data section", comment)
	fmt.Println(eofObject.Data)
}
