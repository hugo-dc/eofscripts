package common

import (
	"context"
	"math/big"
	"strconv"
	"strings"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO: Change name to GetHexBytes
func GetBytes(data string) []string {
	data = strings.Trim(data, " ")
	if data[:2] == "0x" {
		data = data[2:]
	}
	totalBytes := len(data) / 2
	hexBytes := make([]string, totalBytes)

	i := 0
	j := 0
	for i < len(data) {
		hexBytes[j] = string(data[i]) + string(data[i+1])
		i += 2
		j += 1
	}

	return hexBytes
}

func IntToHex(number int64) string {
	hex := strconv.FormatInt(number, 16)

	if len(hex)%2 != 0 {
		hex = "0" + hex
	}
	return hex
}

func SendTransaction(to ethCommon.Address, amount big.Int, data []byte) error {
	ctx := context.Background()
	chainId := int64(1231209)
	gasLimit := uint64(2000000)
	gasPrice := big.NewInt(10 * 18000)

	addr := ethCommon.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	privateKey, err := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")

	cl, err := ethclient.Dial("http://192.168.10.243:8545")
	nonce, err := cl.NonceAt(ctx, addr, nil)

	if err != nil {
		return err
	}
	signer := types.NewEIP2930Signer(big.NewInt(chainId))

	txData := &types.AccessListTx{
		ChainID:  big.NewInt(chainId),
		Nonce:    nonce,
		To:       &to,
		Value:    &amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}

	signedTx, err := types.SignNewTx(privateKey, signer, txData)

	if err != nil {
		return err
	}

	err = cl.SendTransaction(ctx, signedTx)

	if err != nil {
		return err
	}

	return nil
}
