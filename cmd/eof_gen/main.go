package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

func showUsage() {
	fmt.Println("eof_gen - Generate EOF version of the provided EVM code")
	fmt.Println("Usage:")
	fmt.Println("\teof_gen d:[data] c:<code>|C:<input>:<outputs>:<code>")
}

func main() {
	data := ""
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
			bytecode, err := hex.DecodeString(code)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			eofObject.AddCode(bytecode)
		}

		if arg[:2] == "C:" {
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
				bytecode, err := hex.DecodeString(code)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				eofObject.AddCodeWithType(bytecode, code_type)
			}
		}
		if arg[:2] == "K:" {
			container := arg[2:]
			containerBytecode, err := hex.DecodeString(container)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			eofObject.AddContainer(containerBytecode)
		}

	}

	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	eofObject.AddData(dataBytes)
	eof_code := eofObject.Code()
	fmt.Println(eof_code)
}
