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

	opCodes := common.GetOpcodesByName()

	// Push data length
	pushLen := len(data_len_hex) / 2
	pushOp := fmt.Sprintf("PUSH%d", pushLen)

	code_len_hex := common.IntToHex(int64(len(code) + 10 + (pushLen * 2) + 19)) // adds 10 counting the following opcodes and the EOF header
	code = append(code, opCodes[pushOp].AsHex())
	code = append(code, data_len_hex)

	// Push data offset (eof_header_len + evm_bytecode_len)
	code = append(code, opCodes["PUSH1"].AsHex())
	code = append(code, code_len_hex)

	// Push result offset (0)
	code = append(code, opCodes["PUSH1"].AsHex())
	code = append(code, "00")

	// codecopy
	code = append(code, opCodes["CODECOPY"].AsHex())

	// Push return size
	code = append(code, opCodes[pushOp].AsHex())
	code = append(code, data_len_hex)

	// Push return offset (0)
	code = append(code, opCodes["PUSH1"].AsHex())
	code = append(code, "00")

	// return
	code = append(code, opCodes["RETURN"].AsHex())

	data_content := strings.Join(data[:], "")
	code_contents := []string{strings.Join(code[:], "")}

	eofObject := common.NewEOFObject()
	eofObject.AddData(data_content)
	for _, c := range code_contents {
		eofObject.AddCode(c)
	}

	fmt.Println(eofObject.Code())
}
