package common

import (
	"fmt"
	"strconv"
)

// TODO: This will be deleted
func GenerateEOF(data string, types [][]int64, code []string, withTypes bool) string {
	if len(data)%2 != 0 {
		panic("Error: odd data size")
	}

	if len(types) > 1 {
		withTypes = true
	}

	types_length := len(types) * 2
	types_length_hex := strconv.FormatInt(int64(types_length), 16)

	for {
		if len(types_length_hex) != 4 {
			types_length_hex = "0" + types_length_hex
		} else {
			break
		}
	}
	types_section := "03" + types_length_hex

	types_content := ""
	for _, t := range types {
		inputs_hex := strconv.FormatInt(t[0], 16)
		if len(inputs_hex)%2 != 0 {
			inputs_hex = "0" + inputs_hex
		}

		outputs_hex := strconv.FormatInt(t[1], 16)
		if len(outputs_hex)%2 != 0 {
			outputs_hex = "0" + outputs_hex
		}

		types_content = types_content + inputs_hex + outputs_hex
	}

	code_section := ""
	code_content := ""
	for _, c := range code {
		if len(c)%2 != 0 {
			panic("Error: odd code size")
		}

		code_len := len(c) / 2
		code_len_hex := fmt.Sprintf("%x", code_len)
		if len(code_len_hex)%2 != 0 {
			code_len_hex = "0" + code_len_hex
		}

		if len(code_len_hex) < 4 {
			code_len_hex = "00" + code_len_hex
		}
		code_section = code_section + "01" + code_len_hex
		code_content = code_content + c
	}

	data_len := len(data) / 2
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

	terminator := "00"

	eof_code := "ef0001"
	if withTypes {
		eof_code = eof_code + types_section + code_section + data_section + terminator + types_content + code_content + data
	} else {
		eof_code = eof_code + code_section + data_section + terminator + code_content + data
	}

	return eof_code
}
