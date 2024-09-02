package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hugo-dc/eofscripts/common"
)

func main() {
	eof_code := ""
	//is_initcode := false
	if len(os.Args) < 2 {
		fmt.Scanln(&eof_code)
	} else {
		stdin_input := ""
		fmt.Scanln(&stdin_input)
		arg1 := os.Args[1]
		if len(stdin_input) > 1 && arg1 != "--initcode" {
			fmt.Println("Error: No arguments allowed when reading from stdin")
			return
		}
		if arg1 == "--initcode" {
			//is_initcode = true
			eof_code = stdin_input
		} else {
			eof_code = strings.Join(os.Args[1:], " ")
		}
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
	//eofObject.Initcode = is_initcode

	description := eofObject.DescribeAsPython()
	fmt.Println(description)

	/*
			// Print the EOF Object as Python format
			pyeof_code := `Container(
		  name = 'EOFV0001',
		  sections = [
		    %s%s%s  ],
		  kind=ContainerKind.%s
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

					formatted_subcontainer, error := common.ParseEOF(v)

					if error != nil {
						fmt.Println(">> Error: ", error)
						container_sections += fmt.Sprintf(`  Section.Container(
		          container=Container(
		              name="EOFV1_SUBCONTAINER_%v",
		              raw_bytes=bytes(
		                  [ %s ])
		          )
		      ),
		    `, i+0, raw_bytecode)
					} else {
						container_sections += fmt.Sprintf(`  Section.Container(
						  %s
					),
				`, formatted_subcontainer)
					}
				}
			}

			data_section := ""
			if len(eofObject.Data) > 0 {
				data_section += fmt.Sprintf("  Section.Data(data=\"%s\")\n", eofObject.Data)
			}

			fmt.Printf(pyeof_code, code_sections, container_sections, data_section, container_kind)
	*/

}
