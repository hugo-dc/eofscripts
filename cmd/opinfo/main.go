package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("opinfo <opcode> [option]")
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

	if len(opId) > 2 {
		if opId[:2] == "0x" {
			opcodes := common.GetOpcodesByNumber()

			opId64, err := strconv.ParseInt(opId[2:], 16, 64)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if op, ok := opcodes[int(opId64)]; ok {
				opcode = op
			} else {
				fmt.Println("Opcode", opId, "not found!")
				return
			}
		} else {
			opcodes := common.GetOpcodesByName()

			if op, ok := opcodes[opId]; ok {
				opcode = op
			}
		}
	} else {
		showUsage()
		return
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
}
