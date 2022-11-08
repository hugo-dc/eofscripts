package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("returneofdata - Receives data as bytes, will generate EOF code to read the data and return it")
	fmt.Println("Usage:")
	fmt.Println("\treturneofdata <data>")
}

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	data_input := os.Args[1]
	data := common.GetBytes(data_input)
	data_len := len(data)
	data_len_hex := strconv.FormatInt(int64(data_len), 16)

	if len(data_len_hex) < 2 {
		data_len_hex = "0" + data_len_hex
	}

	var result []string

	// Push data length
	result = append(result, common.Push1().AsHex())
	result = append(result, data_len_hex)

	// Push data offset (eof_header_len(10) + evm_bytecode_len(12) = 22 - 0x16)
	result = append(result, common.Push1().AsHex())
	result = append(result, "16")

	// Push result offset (0)
	result = append(result, common.Push1().AsHex())
	result = append(result, "00")

	// codecopy
	result = append(result, common.CodeCopy().AsHex())

	// Push return size
	result = append(result, common.Push1().AsHex())
	result = append(result, data_len_hex)

	// Push return offset (0)
	result = append(result, common.Push1().AsHex())
	result = append(result, "00")

	// return
	result = append(result, common.Return().AsHex())

	eof_code := common.GenerateEOF(strings.Join(data[:], ""), strings.Join(result[:], ""))

	fmt.Println(eof_code)
}
