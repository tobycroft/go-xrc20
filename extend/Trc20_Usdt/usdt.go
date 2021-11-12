package Trc20_Usdt

import (
	"bytes"
	"fmt"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"main.go/config/app_conf"
	"main.go/extend/trx-sign-go-1.0.3/grpcs"
	"main.go/extend/trx-sign-go-1.0.3/sign"
	"math/big"
)

type TokenTransaction struct {
	client          *grpcs.Client
	contractAddress string
}

func InitTranns(contractAddress string) *TokenTransaction {
	//EthRPC_API := SystemParamModel.Api_find_val("EthRPC_API").(string)
	TrcRPC_API := app_conf.TrcRPC_API
	c, err := grpcs.NewClient(TrcRPC_API)
	if err != nil {
		panic(err)
	}
	return &TokenTransaction{client: c, contractAddress: contractAddress}
}

func (c *TokenTransaction) TransferFrom(from, to, contract string, amount *big.Int, feeLimit int64) error {
	a, err := ethabi.JSON(bytes.NewReader([]byte(abiJson)))
	if err != nil {
		fmt.Println("JSON", err)
		return err
	}
	//method:=a.Methods["transferFrom"]
	fromaddress, err := addr.Base58ToAddress(from)
	if err != nil {
		fmt.Println(err)
		return err
	}
	toaddress, err := addr.Base58ToAddress(to)
	if err != nil {
		fmt.Println(err)
		return err
	}
	amount = amount.Mul(amount, big.NewInt(1000000))
	bz, err := a.Pack("transferFrom", common.BytesToAddress(fromaddress.Bytes()), common.BytesToAddress(toaddress.Bytes()), amount)
	//bz, err := abi2.Pack("transferFrom", ab)
	if err != nil {
		fmt.Println("Pack", err)
		return err
	}
	s := common.Bytes2Hex(bz)

	//amount := big.NewInt(20)
	amount = amount.Mul(amount, big.NewInt(1000000))
	tx, err := c.client.GRPC.TRC20Call(from, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", s, false, 2000000)
	if err != nil {
		fmt.Println("TRC20Call", err)
		return err
	}
	//tx, err := c.TransferTrc20("TFTGMfp7hvDtt4fj3vmWnbYsPSmw5EU8oX", "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT", "TJ93jQZibdB3sriHYb5nNwjgkPPAcFR7ty", amount, 100000000)
	signTx, err := sign.SignTransaction(tx.Transaction, "5c023564aa0c582e9a5d127133e9b45c5b9a7a409b22f7e8a5c19d4d3f424eea")
	if err != nil {
		fmt.Println("SignTransaction", err)
		return err
	}
	//fmt.Println("signTx", signTx.String())
	err = c.client.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Println("BroadcastTransaction", err)
		return err
	}
	fmt.Println(common.Bytes2Hex(tx.GetTxid()))

	return nil

}
