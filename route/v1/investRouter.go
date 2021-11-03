package v1

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/exchange/controller"
)

func InvestRouter(route *gin.RouterGroup) {
	route.Any("/", func(context *gin.Context) {
		context.String(0, route.BasePath())
	})

	controller.PaymentController(route.Group("payment"))
	controller.RecordController(route.Group("record"))
}
