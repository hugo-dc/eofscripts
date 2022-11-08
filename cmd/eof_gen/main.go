package main

import (
	"fmt"
	"os"

	common "github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("eof_gen - Generate EOF version of the provided EVM code")
	fmt.Println("Usage:")
	fmt.Println("\teof_gen <data> <code>")
}

func main() {
	if len(os.Args) != 3 {
		showUsage()
		return
	}

	data := os.Args[1]
	code := os.Args[2]

	eof_code := common.GenerateEOF(data, code)

	fmt.Println(eof_code)
}
