package main

import (
	"fmt"
	"main.go/extend/Trc20_Usdt"
	"main.go/tuuz/Calc"
)

func main() {

	//InvestTransfer.InvestTransfer_trc()

	//mainroute := gin.Default()
	//route.OnRoute(mainroute)
	//mainroute.Run(":80")

	t := Trc20_Usdt.InitTranns("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	err, txs := t.TransferFrom("c2e34562e0478a3e4e8f1f79f0d9f156c81249da3df00013531191888a18d7cf", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Calc.ToDecimal(1000000))
	fmt.Println(err, txs)

}
