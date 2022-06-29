package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		showUsage()
		return
	}

	opN, _ := strconv.Atoi(os.Args[1])

	if opN < 1 || opN > 32 {
		showUsage()
		return
	}

	op := 95 + opN
	termOps := [5]string{"00", "f3", "fd", "fe", "ff"}

	//fmt.Println("----------------------------------")
	//fmt.Println(fmt.Sprintf("PUSH%d", opN))
	pushOp := fmt.Sprintf("%x", op)
	fmt.Println(pushOp)
	fmt.Println(fmt.Sprintf("push%d(", opN))

	pushData := ""
	k := 0
	for j := opN; j > 1; j-- {
		pushData += termOps[k]
		k += 1
		if k == 5 {
			k = 0
		}
	}
	pushStatement := pushOp + pushData
	fmt.Println(pushStatement)
	if len(pushData) >= 2 {
		fmt.Println(fmt.Sprintf("push%d", opN) + "(0x" + pushData)
	} else {
		fmt.Println(fmt.Sprintf("push%d", opN) + "(")
	}

	op += 1
}

func showUsage() {
	fmt.Println("gentruncpush - Generate code for truncated PUSH opcodes")
	fmt.Println("Usage:")
	fmt.Println("\tgentruncpush <N>")
	fmt.Println("\t\tWhere N is in pushN (1-32)")
}
