package Trc20_Usdt

import (
	"fmt"
	abi2 "github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"main.go/extend/trx-sign-go-1.0.3/grpcs"
	"main.go/extend/trx-sign-go-1.0.3/sign"
	"math/big"
	"testing"
)

func TransferFrom(t *testing.T) {
	abi, err := abi2.LoadFromJSON(abiJson)
	if err != nil {
		fmt.Println(err)
		return
	}
	bz, err := abi2.Pack("transferFrom", abi)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := common.Bytes2Hex(bz)

	addrB, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, err
	}
	ab := common.LeftPadBytes(amount.Bytes(), 32)
	req := trc20TransferMethodSignature + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrB.Hex())-4:] + addrB.Hex()[4:]
	req += common.Bytes2Hex(ab)

	c, err := grpcs.NewClient("54.168.218.95:50051")
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(20)
	amount = amount.Mul(amount, big.NewInt(1000000000000000000))
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
