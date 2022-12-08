package main

import (
	"fmt"
	"os"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("mnem2evm - Converts Opcode mnemonics to EVM bytecode")
	fmt.Println("Usage:")
	fmt.Println("\tmnem2evm <mnemonics>")
}

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	mnems := os.Args[1]
	evm := common.Mnem2Evm(mnems)
	fmt.Println(evm)
}
