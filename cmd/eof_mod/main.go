package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/eofscripts/common"
)

func showUsage() {
	fmt.Println("eof_mod - Modify provided eof code, adding new sections")
	fmt.Println("Usage:")
	fmt.Println("\teof_mod <eof_code> d:[data] c:<code> C:<inputs><outputs>:<code> t:<types>")
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
		return
	}

	eofCode := os.Args[1]
	eofBytecode, err := hex.DecodeString(eofCode)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	eofObject, err := common.ParseEOF(eofBytecode)

	if err != nil {
		fmt.Println("Err:", err)
	}

	for i := 2; i < len(os.Args); i++ {
		if os.Args[i][:2] == "d:" {
			if len(eofObject.Data) != 0 {
				fmt.Println("Warning: EOF Object alredy contains data, data will be overwritten")
			}
			dataStr := os.Args[i][2:]
			data, err := hex.DecodeString(dataStr)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			eofObject.AddData(data)
		}
		if os.Args[i][:2] == "c:" {
			codeStr := os.Args[i][2:]
			code, err := hex.DecodeString(codeStr)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			eofObject.AddCode(code)
		}
		if os.Args[i][:2] == "C:" {
		}
		if os.Args[i] == ":t" {
			ts := os.Args[i][2:]
			if len(ts)%8 != 0 {
				fmt.Println("Wrong types")
			}
			for len(ts) > 0 {
				inputs, err := strconv.ParseInt(ts[:2], 16, 64)
				if err != nil {
					fmt.Println("Print invalid inputs in type")
				}
				outputs, err := strconv.ParseInt(ts[2:4], 16, 64)
				if err != nil {
					fmt.Println("Print invalid outputs in type")
				}

				maxStack, err := strconv.ParseInt(ts[4:8], 16, 64)
				if err != nil {
					fmt.Println("Print invalid max stack in type")
				}

				eofObject.Types = append(eofObject.Types, []int64{inputs, outputs, maxStack})

			}
		}
	}
	fmt.Println(eofObject.Code())
}
