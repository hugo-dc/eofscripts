package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/ethscripts/common"
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
	eofObject, err := common.ParseEOF(eofCode)

	if err != nil {
		fmt.Println("Err:", err)
	}

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
	fmt.Println(eofObject.Code(false, true))
}
