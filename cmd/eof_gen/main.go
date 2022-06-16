package main

import (
	"fmt"
	"os"
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

	if len(data)%2 != 0 {
		panic("Error: odd data size")
	}

	if len(code)%2 != 0 {
		panic("Error: odd code size")
	}

	code_len := len(code) / 2
	data_len := len(data) / 2

	code_len_hex := fmt.Sprintf("%x", code_len)
	data_len_hex := fmt.Sprintf("%x", data_len)

	if len(code_len_hex)%2 != 0 {
		code_len_hex = "0" + code_len_hex
	}

	if len(code_len_hex) < 4 {
		code_len_hex = "00" + code_len_hex
	}

	fmt.Println("Code length", code_len, code_len_hex)
	code_section := "01" + code_len_hex
	data_section := ""

	if data_len > 0 {
		data_len_hex := fmt.Sprintf("%x", data_len)

		if len(data_len_hex)%2 != 0 {
			data_len_hex = "0" + data_len_hex
		}

		if len(data_len_hex) < 4 {
			data_len_hex = "00" + data_len_hex
		}

		data_section = "02" + data_len_hex
	}

	fmt.Println("Data length", data_len, data_len_hex)

	terminator := "00"

	eof_code := "ef0001" + code_section + data_section + terminator + code + data

	fmt.Println(eof_code)
}
