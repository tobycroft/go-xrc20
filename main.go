package main

import (
	"main.go/extend/Trc20_Usdt"
	"math/big"
)

func main() {

	//InvestTransfer.InvestTransfer_trc()
t:=Trc20_Usdt.InitTranns("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
t.TransferFrom()
	Trc20_Usdt.TransferFrom("TBA6CypYJizwA9XdC7Ubgc5F1bxrQ7SqPt", "TLFtevwBqEV9fkC4nNwmTRtFiygPJKpqci", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", big.NewInt(22884), 0)
	//mainroute := gin.Default()
	//route.OnRoute(mainroute)
	//mainroute.Run(":80")

}
