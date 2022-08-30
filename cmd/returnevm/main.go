package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("returnevm - Wraps EVM bytecode to be returned by a deployer contract")
	fmt.Println("Usage:")
	fmt.Println("\treturnevm <evm bytecode>")
}

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	data := os.Args[1]
	code := common.GetBytes(data)

	codelen := len(code)

	pushlen := codelen
	if codelen > 32 {
		pushlen = 32
	}

	codelenhex := strconv.FormatInt(int64(codelen), 16)
	if len(codelenhex)%2 != 0 {
		codelenhex = "0" + codelenhex
	}

	pushOp := strconv.FormatInt(int64(95+pushlen), 16)

	var result []string
	result = append(result, pushOp)                  // PUSHn
	result = append(result, code[:]...)              // Code    (Value)
	result = append(result, common.Push1().AsHex())  // PUSH1
	result = append(result, "00")                    // 00 (Offset)
	result = append(result, common.MStore().AsHex()) // MSTORE
	result = append(result, common.Push1().AsHex())  // PUSH1
	result = append(result, codelenhex)              // CodeLength (Offset end)
	result = append(result, common.Push1().AsHex())  // PUSH1

	if codelen < 32 {
		initialOffset := strconv.FormatInt(int64(32-codelen), 16)

		if len(initialOffset)%2 != 0 {
			initialOffset = "0" + initialOffset
		}
		result = append(result, initialOffset) // Offset
	} else {
		result = append(result, "00") // Offset
	}
	result = append(result, common.Return().AsHex()) // RETURN
	result = append(result, common.Stop().AsHex())   // STOP

	// Show result
	fmt.Println(strings.Join(result[:], ""))
}
