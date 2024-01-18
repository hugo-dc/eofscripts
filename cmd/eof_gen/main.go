package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("eof_gen - Generate EOF version of the provided EVM code")
	fmt.Println("Usage:")
	fmt.Println("\teof_gen d:[data] c:<code>|C:<input>:<outputs>:<code> [-t]")
}

func main() {
	data := ""
	showTypes := false
	oldCode := false
	eofObject := common.NewEOFObject()
	for _, arg := range os.Args {
		if arg[:2] == "d:" {
			data = arg[2:]
		}

		if arg[:2] == "c:" {
			code := arg[2:]
			_, err := hex.DecodeString(code)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			eofObject.AddCode(code)
		}

		if arg[:2] == "C:" {
			showTypes = true
			code_contents := strings.Split(arg, ":")

			if len(code_contents) != 4 {
				fmt.Println("Error: Expect typed code C:<input>,<outputs>,<code>")
				return
			} else {
				inputs, err := strconv.ParseInt(code_contents[1], 10, 64)

				if err != nil {
					fmt.Println("Error: ", err)
					return
				}

				outputs, err := strconv.ParseInt(code_contents[2], 10, 64)

				if err != nil {
					fmt.Println("Error: ", err)
					return
				}

				code_type := []int64{inputs, outputs}
				code := code_contents[3]
				_, err = hex.DecodeString(code)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				eofObject.AddCodeWithType(code, code_type)
			}
		}

		if arg == "-t" {
			showTypes = true
		}

		if arg == "-o" {
			oldCode = true
		}
	}

	if !oldCode {
		showTypes = true
	}

	eofObject.AddData(data)
	eof_code := eofObject.Code(oldCode, showTypes)
	fmt.Println(eof_code)
}
