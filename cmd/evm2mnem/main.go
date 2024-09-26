package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/hugo-dc/eofscripts/common"
)

func showUsage() {
	fmt.Println("evm2mnem - Converts EVM Bytecode to mnemonics")
	fmt.Println("Usage:")
	fmt.Println("\tevm2mnem <bytecode>")
}

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	bytecodeStr := os.Args[1]
	if bytecodeStr[:2] == "0x" {
		bytecodeStr = bytecodeStr[2:]
	}

	bytecode, err := hex.DecodeString(bytecodeStr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	opcalls, err := common.BytecodeToOpCalls(bytecode)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		mnem := common.Evm2Mnem(opcalls)
		fmt.Println(mnem)
	}
}
