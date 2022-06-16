package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	ethCommon "github.com/ethereum/go-ethereum/common"
	common "github.com/hugo-dc/ethscripts/common"
)

func main() {
	if len(os.Args) != 4 {
		showUsage()
		return
	}

	to := ethCommon.HexToAddress(os.Args[1])
	amount, err := strconv.ParseInt(os.Args[2], 10, 64)

	if err != nil {
		fmt.Println("`amount` is not valid ")
	} else {
		amount := big.NewInt(amount)
		data := ethCommon.Hex2Bytes(os.Args[3])

		err := common.SendTransaction(to, *amount, data)

		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	}

}

func showUsage() {
	fmt.Println("call - Call address")
	fmt.Println("Usage:")
	fmt.Println("\t call <address> <amount> <data>")
}
