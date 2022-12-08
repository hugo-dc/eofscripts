package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("returneofdata - Receives data as bytes, will generate EOF code to read the data and return it, if some code is provided, it will be added before the code returning the data")
	fmt.Println("Usage:")
	fmt.Println("\treturneofdata <data> [code]")
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		showUsage()
		return
	}

	data_input := os.Args[1]
	data := common.GetBytes(data_input)
	data_len_hex := common.IntToHex(int64(len(data)))

	var code []string
	if len(os.Args) == 3 {
		code_input := os.Args[2]
		code = common.GetBytes(code_input)
	}

	code_len_hex := common.IntToHex(int64(len(code) + 22)) // adds 12 counting the following opcodes and the EOF header

	// Push data length
	code = append(code, common.Push1().AsHex())
	code = append(code, data_len_hex)

	// Push data offset (eof_header_len + evm_bytecode_len)
	code = append(code, common.Push1().AsHex())
	code = append(code, code_len_hex)

	// Push result offset (0)
	code = append(code, common.Push1().AsHex())
	code = append(code, "00")

	// codecopy
	code = append(code, common.CodeCopy().AsHex())

	// Push return size
	code = append(code, common.Push1().AsHex())
	code = append(code, data_len_hex)

	// Push return offset (0)
	code = append(code, common.Push1().AsHex())
	code = append(code, "00")

	// return
	code = append(code, common.Return().AsHex())

	data_content := strings.Join(data[:], "")
	code_contents := []string{strings.Join(code[:], "")}
	types := [][]int64{}

	eof_code := common.GenerateEOF(data_content, types, code_contents, false)
	fmt.Println(eof_code)
}
