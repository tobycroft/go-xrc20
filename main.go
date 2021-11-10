package main

import (
	"github.com/gin-gonic/gin"
	"main.go/route"
)

func main() {

	//InvestTransfer.InvestTransfer_trc()

	mainroute := gin.Default()
	route.OnRoute(mainroute)
	mainroute.Run(":80")

}
