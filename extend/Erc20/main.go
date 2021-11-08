package Erc20

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"io/ioutil"
	"main.go/common/BaseModel/SystemParamModel"
	"math"
	"math/big"
	"strings"
)

type TokenTransaction struct {
	client          *ethclient.Client
	contractAddress string
}

func InitTranns(contractAddress string) *TokenTransaction {
	EthRPC_API := SystemParamModel.Api_find_val("EthRPC_API")
	rpcDial, err := rpc.Dial(EthRPC_API.(string))
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(rpcDial)
	return &TokenTransaction{client: client, contractAddress: contractAddress}
}

func (s *TokenTransaction) Transaction(toAddress, keyfile, pwd string, tokenAmount float64) (err error) {
	i, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(i)), pwd)
	if err != nil {
		return
	}

	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return
	}

	amount := big.NewFloat(tokenAmount)
	tenDecimal := big.NewFloat(math.Pow(10, float64(18)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, amount).Int(&big.Int{})
	auth.GasLimit = 200000
	txs, err := token.Transfer(auth, common.HexToAddress(toAddress), convertAmount)
	if err != nil {
		return
	}

	fmt.Println("chainId---->", txs.ChainId())
	return
}
