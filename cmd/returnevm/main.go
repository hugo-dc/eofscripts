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
	codelenhex := strconv.FormatInt(int64(codelen), 16)
	if len(codelenhex)%2 != 0 {
		codelenhex = "0" + codelenhex
	}

	// If code length is greater than 32, calculate how many chunks of 32 bytes
	// each can be created
	totalChunks := 1
	if codelen > 32 {
		totalChunks = codelen / 32

		if codelen%32 != 0 {
			totalChunks += 1
		}
	}

	// Split code into chunks
	codeChunks := [][]string{}
	for i := 0; i < totalChunks; i++ {
		start := (i * 32)
		end := (i*32 + 32)

		if end > codelen {
			end = codelen
		}

		chunk := code[start:end]

		codeChunks = append(codeChunks, chunk)
	}

	// Store code in memory, chunk by chunk
	opCodes := common.GetOpcodesByName()
	var result []string
	for i := 0; i < len(codeChunks); i++ {
		chunk := codeChunks[i]
		pushlen := len(chunk)

		pushOp := strconv.FormatInt(int64(95+pushlen), 16)
		result = append(result, pushOp)      // PUSHn
		result = append(result, chunk[:]...) // Code

		if pushlen < 32 {
			totalBits := (32 - pushlen) * 8

			totalBitsHex := strconv.FormatInt(int64(totalBits), 16)
			if len(totalBitsHex)%2 != 0 {
				totalBitsHex = "0" + totalBitsHex
			}
			result = append(result, opCodes["PUSH1"].AsHex())
			result = append(result, totalBitsHex) // <totalBits>
			result = append(result, opCodes["SHL"].AsHex())
		}

		result = append(result, opCodes["PUSH1"].AsHex())

		offset := strconv.FormatInt(int64(i*32), 16)
		if len(offset)%2 != 0 {
			offset = "0" + offset
		}
		result = append(result, offset) // Offset
		result = append(result, opCodes["MSTORE"].AsHex())
	}

	result = append(result, opCodes["PUSH1"].AsHex())
	result = append(result, codelenhex) // CodeLength (Offset end)
	result = append(result, opCodes["PUSH1"].AsHex())

	result = append(result, "00") // Offset
	result = append(result, opCodes["RETURN"].AsHex())
	result = append(result, opCodes["STOP"].AsHex())

	// Show result
	fmt.Println(strings.Join(result[:], ""))
}
