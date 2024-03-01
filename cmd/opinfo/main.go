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
	fmt.Println("    <opcode>\t\tOpcode name, decimal, or hexadecimal value (0x prefixed)")
	fmt.Println("  options:")
	fmt.Println("    --name\t\tShow opcode Name")
	fmt.Println("    --hex\t\tShow opcode hexadecimal value")
	fmt.Println("    --inputs\t\tShow required stack items")
	fmt.Println("    --outputs\t\tShow pushed stack items")
	fmt.Println("    --is-push\t\tReturns true if opcode is a PUSH opcode")
	fmt.Println("    --immediates\tReturns the number of bytes the opcodes requires as immediates")
	fmt.Println("    --is-terminating\tReturns true if opcode is a terminating opcode")
	fmt.Println("    --all\t\tShow all information")
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
	switch option {
	case "--name":
		fmt.Println(opcode.Name)
	case "--hex":
		fmt.Printf("%02x\n", opcode.Code)
	case "--inputs":
		fmt.Println(opcode.StackInput)
	case "--outputs":
		fmt.Println(opcode.StackOutput)
	case "--is-push":
		if opcode.IsPush() {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	case "--immediates":
		fmt.Println(opcode.Immediates)
	case "--is-terminating":
		if opcode.IsTerminating {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	case "--all":
		fmt.Println("Name:", opcode.Name)
		fmt.Println("Hex:", fmt.Sprintf("0x%02x", opcode.Code))
		fmt.Println("Inputs:", opcode.StackInput)
		fmt.Println("Outputs:", opcode.StackOutput)
		fmt.Println("IsPush:", opcode.IsPush())
		fmt.Println("Immediates:", opcode.Immediates)
		fmt.Println("IsTerminating:", opcode.IsTerminating)
	default:
		showUsage()
	}
}
