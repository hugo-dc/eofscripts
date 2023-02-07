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
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(10 * 18000)

	addr := ethCommon.HexToAddress("0x0a0e72E7Ec1e636Cd90C723E4d4b249dc7A08d37")
	privateKey, err := crypto.HexToECDSA("a78ba7ce1fac579d25e98b3b9b80382e0684dff695bd61150abdc6726677b510")

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
