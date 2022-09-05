package main

import (
	"fmt"
	"os"

	ethCommon "github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

func main() {
	var salt [32]byte

	if len(os.Args) != 4 {
		showUsage()
		return
	}

	originAddr := ethCommon.HexToAddress(os.Args[1])
	salt_tmp := ethCommon.FromHex(os.Args[2])

	for i := 0; i < len(salt_tmp); i++ {
		salt[32-len(salt_tmp)+i] = salt_tmp[i]
	}

	code := ethCommon.FromHex(os.Args[3])
	codeHash := crypto.Keccak256(code)
	address := crypto.CreateAddress2(originAddr, salt, codeHash)

	fmt.Println(address)

}

func showUsage() {
	fmt.Println("create_address - returns the new address created by an origin given a nonce")
	fmt.Println("Usage:")
	fmt.Println("\t create_address <origin_address> <salt (hex)> <code (hex)>")
}
