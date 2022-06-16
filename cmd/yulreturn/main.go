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

	code := os.Args[1]

	if len(code)%2 != 0 {
		panic("Error: Code length is even")
	}

	code_len := len(code) / 2

	yul_code := code

	for i := code_len; i < 32; i++ {
		yul_code = yul_code + "00"
	}

	yul_code = "{ mstore(0, 0x" + yul_code + ") return (0, " + fmt.Sprintf("%d", code_len) + ") }"

	fmt.Println(yul_code)

}

func showUsage() {
	fmt.Println("yulreturn - Receives evm bytecode as parameter, returns yul code which returns the evm bytecode")
	fmt.Println("Usage:")
	fmt.Println("\tyulreturn <evm_bytecode>")
}
