package main

import (
	"fmt"
	"os"

	"github.com/hugo-dc/ethscripts/common"
)

func usage() {
	fmt.Println("eof_upd - Update EOF code to match current unified spec header")
	fmt.Println("\tUsage:")
	fmt.Println("\teof_upd <EOF_Code>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	eofCode := os.Args[1]
	eofObject := common.ParseOldEOF(eofCode)
	fmt.Println(eofObject.CodeNew(false))
}
