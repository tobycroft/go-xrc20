package Erc20

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	EthRPC_API := app_conf.EthRPC_API
	rpcDial, err := rpc.Dial(EthRPC_API)
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(rpcDial)
	return &TokenTransaction{client: client, contractAddress: contractAddress}
}

func (s *TokenTransaction) Transaction(privateKey string, fromAddress, toAddress string, tokenAmount float64) (err error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return fmt.Errorf("hex decode private key error: %v", err)
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	//auth, err := bind.NewTransactor(strings.NewReader(string(i)), pwd)
	auth := bind.NewKeyedTransactor(priv)
	//if err != nil {
	//	return
	//}

	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return
	}

	amount := big.NewFloat(tokenAmount)
	tenDecimal := big.NewFloat(math.Pow(10, float64(18)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, amount).Int(&big.Int{})
	auth.GasLimit = 200000
	//txs, err := token.Transfer(auth, common.HexToAddress(toAddress), convertAmount)
	//if err != nil {
	//	return
	//}
	txs, err := token.TransferFrom(auth, common.HexToAddress(fromAddress), common.HexToAddress(toAddress), convertAmount)
	if err != nil {
		return
	}
	fmt.Println("chainId---->", txs.ChainId())
	return
}
