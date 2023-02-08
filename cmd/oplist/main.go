package main

import (
	"fmt"
	"sort"

	"github.com/hugo-dc/ethscripts/common"
)

func main() {
	opcodes := common.GetOpcodesByNumber()

	keys := make([]int, 0)
	for k, _ := range opcodes {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		op := opcodes[k]
		fmt.Printf("0x%02x\n", op.Code)
	}
}
