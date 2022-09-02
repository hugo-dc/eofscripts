package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	ethCommon "github.com/ethereum/go-ethereum/common"
	common "github.com/hugo-dc/ethscripts/common"
)

func showUsage() {
	fmt.Println("call - Call address")
	fmt.Println("Usage:")
	fmt.Println("\t call <address> [amount] [data] - amount default: 0, data default: empty.")
}

func main() {
	if len(os.Args) != 4 && len(os.Args) != 2 {
		showUsage()
		return
	}

	amount := big.NewInt(0)
	data_str := ""
	if len(os.Args) == 4 {
		data_str = os.Args[3]

		if data_str[:2] == "0x" {
			data_str = data_str[2:]
		}

		amount_, err := strconv.ParseInt(os.Args[2], 10, 64)

		if err != nil {
			fmt.Println("`amount` is not valid")
			return
		}

		amount = big.NewInt(amount_)
	}

	to := ethCommon.HexToAddress(os.Args[1])
	data := ethCommon.Hex2Bytes(data_str)

	err := common.SendTransaction(to, *amount, data)

	if err != nil {
		fmt.Println("ERROR: ", err)
	}

}
