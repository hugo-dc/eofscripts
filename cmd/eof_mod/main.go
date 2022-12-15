package main

import (
	"fmt"
	"os"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("eof_mod - Modify provided eof code, adding new sections")
	fmt.Println("Usage:")
	fmt.Println("\teof_mod <eof_code> d:[data] c:<code> C:<inputs><outputs>:<code> [-t]")
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
		return
	}

	eofCode := os.Args[1]

	eofObject := common.ParseEOF(eofCode)

	defaultTypes := false
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i][:2] == "d:" {
			if eofObject.Data != "" {
				fmt.Println("Warning: EOF Object alredy contains data, data will be overwritten")
			}
			eofObject.AddData(os.Args[i][2:])
		}
		if os.Args[i][:2] == "c:" {
			eofObject.AddCode(os.Args[i][2:])
		}
		if os.Args[i][:2] == "C:" {
		}
		if os.Args[i] == "-t" {
			defaultTypes = true
		}
	}
	if defaultTypes {
		if eofObject.AddDefaultType() == false {
			fmt.Println("Warning: Code already contains types")
		}
	}
	fmt.Println(eofObject.Code())
}
