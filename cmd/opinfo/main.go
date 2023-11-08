package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("opinfo <opcode> [option]")
	fmt.Println("  opcode:")
	fmt.Println("    <opcode>\tOpcode name, decimal, or hexadecimal value (0x prefixed)")
	fmt.Println("  options:")
	fmt.Println("    --name\tShow opcode Name")
	fmt.Println("    --hex\tShow opcode hexadecimal value")
	fmt.Println("    --inputs\tShow required stack items")
	fmt.Println("    --outputs\tShow pushed stack items")
	fmt.Println("    --ispush\tReturns true if opcode is a PUSH opcode")
	fmt.Println("    --immediates\tReturns the number of bytes the opcodes requires as immediates")
	fmt.Println("    --is-terminating\tReturns true if opcode is a terminating opcode")
}

func main() {

	if len(os.Args) != 3 {
		showUsage()
		return
	}

	opId := os.Args[1]

	var opcode common.OpCode

	if opId[:2] == "0x" {

		opId64, err := strconv.ParseInt(opId[2:], 16, 64)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		opcodes := common.GetOpcodesByNumber()

		if op, ok := opcodes[int(opId64)]; ok {
			opcode = op
		} else {
			fmt.Println("Opcode", opId, "not found!")
			return
		}
	} else {
		opId64, err := strconv.ParseInt(opId, 10, 64)

		if err != nil {
			opcodes := common.GetOpcodesByName()

			if op, ok := opcodes[opId]; ok {
				opcode = op
			} else {
				fmt.Println("Opcode", opId, "not found!")
				return
			}
		} else {
			fmt.Println("⚠️ Warning: opcode value", opId, "is being interpreted as decimal\n")
			opcodes := common.GetOpcodesByNumber()

			if op, ok := opcodes[int(opId64)]; ok {
				opcode = op
			}
		}
	}

	option := os.Args[2]
	if option == "--name" {
		fmt.Println(opcode.Name)
	}
	if option == "--hex" {
		fmt.Printf("%02x\n", opcode.Code)
	}
	if option == "--inputs" {
		fmt.Println(opcode.StackInput)
	}
	if option == "--outputs" {
		fmt.Println(opcode.StackOutput)
	}
	if option == "--ispush" {
		if opcode.Code >= 0x60 && opcode.Code <= 0x7f {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	}
	if option == "--immediates" {
		fmt.Println(opcode.Immediates)
	}
	if option == "--is-terminating" {
		if opcode.IsTerminating {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	}
	if option == "--all" {
		fmt.Println("Name:", opcode.Name)
		fmt.Println("Hex:", fmt.Sprintf("0x%02x", opcode.Code))
		fmt.Println("Inputs:", opcode.StackInput)
		fmt.Println("Outputs:", opcode.StackOutput)
		fmt.Println("IsPush:", opcode.IsPush)
		fmt.Println("Immediates:", opcode.Immediates)
		fmt.Println("IsTerminating:", opcode.IsTerminating)
	}
}
