package main

import (
	"github.com/gin-gonic/gin"
	"main.go/app/cron/InvestTransfer"
	"main.go/route"
)

func main() {

	go InvestTransfer.InvestTransfer()

	mainroute := gin.Default()
	route.OnRoute(mainroute)
	mainroute.Run(":80")

}
