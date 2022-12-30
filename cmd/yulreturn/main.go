package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	// TODO: Manage code as bytes:
	//     data := os.Args[1]
	//     code := common.GetBytes(data)
	code := os.Args[1]

	if len(code)%2 != 0 {
		panic("Error: Code length is even")
	}

	if code[:2] == "0x" {
		code = code[2:]
	}

	code_len := len(code) / 2

	totalChunks := 1
	if code_len > 32 {
		totalChunks = code_len / 32

		if code_len%32 != 0 {
			totalChunks += 1
		}
	}

	codeChunks := []string{}
	for i := 0; i < totalChunks; i++ {
		start := (i * 32)
		end := (i*32 + 32)

		if end > code_len {
			end = code_len
		}

		start *= 2
		end *= 2

		chunk := code[start:end]
		codeChunks = append(codeChunks, chunk)
	}

	mstores := ""
	for i := 0; i < len(codeChunks); i++ {
		if i > 0 {
			mstores += fmt.Sprintf(" mstore(%d", i*32)
		} else {
			mstores += fmt.Sprintf("mstore(%d", i*32)
		}

		chunk := codeChunks[i]

		if len(chunk) < 64 {
			for i := len(chunk); i < 64; i += 2 {
				chunk = chunk + "00"
			}
		}
		mstores += ", 0x" + chunk + ")"

	}

	yul_code := "{ " + mstores + " return(0, " + fmt.Sprintf("%d", code_len) + ") }"
	fmt.Println(yul_code)
}

func showUsage() {
	fmt.Println("yulreturn - Receives evm bytecode as parameter, returns yul code which returns the evm bytecode")
	fmt.Println("Usage:")
	fmt.Println("\tyulreturn <evm_bytecode>")
}
