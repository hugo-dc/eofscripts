package main

import (
	"fmt"
	"os"
	"strconv"

	ethCommon "github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if len(os.Args) != 3 {
		showUsage()
		return
	}

	originAddr := ethCommon.HexToAddress(os.Args[1])
	nonce, err := strconv.ParseUint(os.Args[2], 10, 64)

	if err != nil {
		fmt.Println("`nonce` is invalid")
		return
	}

	address := crypto.CreateAddress(originAddr, nonce)

	fmt.Println(address)

}

func showUsage() {
	fmt.Println("create_address - returns the new address created by an origin given a nonce")
	fmt.Println("Usage:")
	fmt.Println("\t create_address <origin_address> <nonce>")
}
