package Trc20_Usdt

import (
	"bytes"
	"fmt"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"main.go/extend/trx-sign-go-1.0.3/grpcs"
	"main.go/extend/trx-sign-go-1.0.3/sign"
	"math/big"
)

func TransferFrom(from, to, contract string, amount *big.Int, feeLimit int64) error {
	a, err := ethabi.JSON(bytes.NewReader([]byte(abiJson)))
	if err != nil {
		return err
	}
	//method:=a.Methods["transferFrom"]
	bz, err := a.Pack("transferFrom", from, to, amount)
	//bz, err := abi2.Pack("transferFrom", ab)
	if err != nil {
		return err
	}
	fmt.Println(bz)
	return

	s := common.Bytes2Hex(bz)

	c, err := grpcs.NewClient("54.168.218.95:50051")
	if err != nil {
		return err
	}
	//amount := big.NewInt(20)
	amount = amount.Mul(amount, big.NewInt(1000000000000000000))
	c.GRPC.TRC20Call(from, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", s, false, 2000000)
	tx, err := c.TransferTrc20("TFTGMfp7hvDtt4fj3vmWnbYsPSmw5EU8oX", "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT",
		"TJ93jQZibdB3sriHYb5nNwjgkPPAcFR7ty", amount, 100000000)
	signTx, err := sign.SignTransaction(tx.Transaction, "5c023564aa0c582e9a5d127133e9b45c5b9a7a409b22f7e8a5c19d4d3f424eea")
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}
