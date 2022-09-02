package main

import (
	"context"
	"fmt"
	"math/big"

	"os"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	if len(os.Args) != 2 {
		showUsage()
		return
	}

	data_str := os.Args[1]
	if data_str[:2] == "0x" {
		data_str = data_str[2:]
	}
	data := ethCommon.Hex2Bytes(data_str)
	fmt.Println("data:", data)

	cl, err := ethclient.Dial("http://192.168.10.243:8545")

	addr := ethCommon.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")

	ctx := context.Background()
	//tx := new(types.Transaction)
	nonce, err := cl.NonceAt(ctx, addr, nil)
	//to := ethCommon.HexToAddress("")
	amount := big.NewInt(0)
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(10 * 18000)

	chainId := int64(1231209)

	//tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	//var privateKey ecdsa.PrivateKey
	//privateKey.D = new(big.Int).SetBytes(ethCommon.Hex2Bytes("a78ba7ce1fac579d25e98b3b9b80382e0684dff695bd61150abdc6726677b510"))
	privateKey, err := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")

	if err != nil {
		fmt.Println("Err: ", err)
	}
	//privateKey := "a78ba7ce1fac579d25e98b3b9b80382e0684dff695bd61150abdc6726677b510"
	/*
		tx, _ := types.NewTx(&types.AccessListTx {
			ChainID: big.NewInt(chainId),
			Nonce: nonce,
			To: &to,
			Value: amount,
			Gas: gasLimit,
			GasPrice: gasPrice,
			Data: data,
		})//.WithSignature(
	*/

	signer := types.NewEIP2930Signer(big.NewInt(chainId))

	txData := &types.AccessListTx{
		ChainID:  big.NewInt(chainId),
		Nonce:    nonce,
		To:       nil,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}

	fmt.Println("PK: ", privateKey)
	fmt.Println("Signer:", signer)
	fmt.Println("TxData:", txData)

	signedTx, err := types.SignNewTx(privateKey, signer, txData)

	//sk := crypto.ToECDSAUnsafe(ethCommon.FromHex(privateKey)) // Sign the transaction
	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(nil), sk)

	if err != nil {
		fmt.Println("err: ", err)
	} else {
		err = cl.SendTransaction(ctx, signedTx)
		fmt.Println("err: ", err)
	}

}

func showUsage() {
	fmt.Println("deploy - deploy the provided evm bytecode")
	fmt.Println("Usage:")
	fmt.Println("\t deploy <evm init code>")
}
