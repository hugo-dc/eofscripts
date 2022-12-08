package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("eof_gen - Generate EOF version of the provided EVM code")
	fmt.Println("Usage:")
	fmt.Println("\teof_gen d:[data] c:<code>|C:<input>:<outputs>:<code>")
}

func main() {
	data := ""
	code := []string{}
	types := [][]int64{}

	showTypes := false
	for _, arg := range os.Args {
		if arg[:2] == "d:" {
			data = arg[2:]
		}

		if arg[:2] == "c:" {
			code_type := []int64{0, 0}
			types = append(types, code_type)
			code = append(code, arg[2:])
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
				types = append(types, code_type)

				code = append(code, code_contents[3])
			}
		}
	}

	eof_code := common.GenerateEOF(data, types, code, showTypes)
	fmt.Println(eof_code)
}
