package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

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
	description := eofObject.Describe()
	fmt.Println(description)
}
