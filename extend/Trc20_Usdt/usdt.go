package Trc20_Usdt

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/shopspring/decimal"
	"main.go/config/app_conf"
	"math"
	"math/big"
)

type TokenTransaction struct {
	client          *ethclient.Client
	contractAddress string
}

func InitTranns(contractAddress string) *TokenTransaction {
	//EthRPC_API := SystemParamModel.Api_find_val("EthRPC_API").(string)
	TrcRPC_API := app_conf.TrcRPC_API

	rpcDial, err := rpc.Dial(TrcRPC_API)
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(rpcDial)
	if err != nil {
		return nil
	}

	return &TokenTransaction{client: client, contractAddress: contractAddress}
}

func (s *TokenTransaction) TransferFrom(privateKey string, fromAddress, toAddress string, tokenAmount decimal.Decimal) (err error, txs *types.Transaction) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err, nil
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	//auth, err := bind.NewTransactor(strings.NewReader(string(i)), pwd)
	auth := bind.NewKeyedTransactor(priv)
	//if err != nil {
	//	return
	//}
	token, err := NewUsdtapi(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		fmt.Println("NewUsdtapi", err)
		return err, nil
	}

	tenDecimal := big.NewFloat(math.Pow(10, float64(6)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, tokenAmount.BigFloat()).Int(&big.Int{})
	auth.GasLimit = 200000
	//txs, err := token.Transfer(auth, common.HexToAddress(toAddress), convertAmount)
	//if err != nil {
	//	return
	//}
	addrA, err := address.Base58ToAddress(fromAddress)
	if err != nil {
		return err, nil
	}
	addrB, err := address.Base58ToAddress(toAddress)
	if err != nil {
		return err, nil
	}

	txs, err = token.TransferFrom(auth, common.BytesToAddress(addrA.Bytes()), common.BytesToAddress(addrB.Bytes()), convertAmount)
	if err != nil {
		fmt.Println("TransferFrom", err)
		return err, nil
	}
	fmt.Println("hash", txs.Hash())
	fmt.Println("type", txs.Type())
	return
}
