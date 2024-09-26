package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

func main() {
	eof_code := ""
	is_initcode := false
	explicit_runtime := false
	if len(os.Args) < 2 {
		fmt.Scanln(&eof_code)
	} else {
		stdin_input := ""
		fmt.Scanln(&stdin_input)
		arg1 := os.Args[1]
		if len(stdin_input) > 1 && arg1 != "--initcode" && arg1 != "--runtime" {
			fmt.Println("Error: No arguments allowed when reading from stdin")
			return
		}
		if arg1 == "--initcode" {
			is_initcode = true
			eof_code = stdin_input
		} else {
			if arg1 == "--runtime" {
				explicit_runtime = true
				eof_code = stdin_input
			} else {
				eof_code = strings.Join(os.Args[1:], " ")
			}
		}
	}
	if eof_code[:2] == "0x" {
		eof_code = eof_code[2:]
	}
	eof_code = strings.Replace(eof_code, " ", "", -1)
	eofBytecode, err := hex.DecodeString(eof_code)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	eofObject, err := common.ParseEOF(eofBytecode)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	if is_initcode {
		eofObject.SetInitcode(is_initcode)
	}
	if explicit_runtime {
		eofObject.SetExplicitRuntime(explicit_runtime)
	}

	description := eofObject.DescribeAsPython(0, 0)
	fmt.Println(description)
}
