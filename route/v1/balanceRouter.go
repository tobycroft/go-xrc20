package v1

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/balance/controller"
)

func BalanceRouter(route *gin.RouterGroup) {
	route.Any("/", func(context *gin.Context) {
		context.String(0, route.BasePath())
	})

	controller.TransferController(route.Group("transfer"))
	controller.ReceiveController(route.Group("receive"))
	controller.RecordController(route.Group("record"))
	controller.RefundController(route.Group("refund"))
}
