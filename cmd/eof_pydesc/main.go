package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

func opEqual(op1 common.OpCall, op2 common.OpCall) bool {
	if op1.OpCode != op2.OpCode {
		return false
	}
	if len(op1.Immediates) != len(op2.Immediates) {
		return false
	}
	for i, imm := range op1.Immediates {
		if imm != op2.Immediates[i] {
			return false
		}
	}
	return true
}

func describeCode(code string) (string, error) {
	ops, err := common.DescribeBytecode(code)
	if err != nil {
		return "", err
	}

	prev_op := common.OpCall{}
	op_count := 1
	descriptions := []string{}
	code_desc := ""
	for _, opc := range ops {
		if opEqual(opc, prev_op) {
			op_count++
			continue
		}
		asm := common.Evm2PyAsm([]common.OpCall{opc})
		if op_count > 1 {
			code_desc += fmt.Sprintf("* %v ", op_count)
			descriptions = append(descriptions, code_desc)
			code_desc = fmt.Sprintf("%v", asm)
			op_count = 1
		} else {
			if code_desc != "" {
				descriptions = append(descriptions, code_desc)
			}
			code_desc = fmt.Sprintf("%v", asm)
		}
		prev_op = opc
	}
	if op_count > 1 {
		code_desc += fmt.Sprintf("* %v ", op_count)
		code_desc = ""
	}
	/*
		if code_desc != "" {
			descriptions = append(descriptions, code_desc)
		}
	*/
	descriptions = append(descriptions, code_desc)
	return strings.Join(descriptions, "+ "), nil
}

func main() {
	eof_code := ""
	if len(os.Args) < 2 {
		fmt.Scanln(&eof_code)
	} else {
		stdin_input := ""
		fmt.Scanln(&stdin_input)
		if len(stdin_input) > 0 {
			fmt.Println("Error: No arguments allowed when reading from stdin")
			return
		}
		eof_code = strings.Join(os.Args[1:], " ")
	}
	if eof_code[:2] == "0x" {
		eof_code = eof_code[2:]
	}
	eof_code = strings.Replace(eof_code, " ", "", -1)
	eofObject, err := common.ParseEOF(eof_code)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Print the EOF Object as Python format
	pyeof_code := `Container(
  name = 'EOFV0001',
  sections = [
    %s%s%s  ],
  kind=ContainerKind.RUNTIME
)
`

	code_sections := ""
	for i, v := range eofObject.CodeSections {
		// Get Max Stack Height
		max_stack_height := eofObject.Types[i][2]
		code := ""
		code_description, err := describeCode(v)
		if err != nil {
			code = v + "\n"
		} else {
			code = code_description
		}
		code_sections += fmt.Sprintf("  Section.Code(code=%s, max_stack_height=%v),\n    ", code, max_stack_height)
	}

	container_sections := ""
	if len(eofObject.ContainerSections) > 0 {
		for i, v := range eofObject.ContainerSections {
			raw_bytecode := ""
			for i := 0; i < len(v); i += 2 {
				separator := ", "
				if i == len(v)-2 {
					separator = ""
				}
				raw_bytecode += fmt.Sprintf("0x%s%s", v[i:i+2], separator)
			}

			container_sections = fmt.Sprintf(`  Section.Container(
          container=Container(
              name="EOFV1_SUBCONTAINER_%v",
              raw_bytes=bytes(
                  [ %s ])
          )
      ])
`, i+0, raw_bytecode)
		}
	}

	data_section := ""
	if len(eofObject.Data) > 0 {
		data_section += fmt.Sprintf("  Section.Data(data=\"%s\")\n", eofObject.Data)
	}

	fmt.Printf(pyeof_code, code_sections, container_sections, data_section)

}
