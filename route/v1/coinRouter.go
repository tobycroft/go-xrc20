package v1

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/coin/controller"
)

func CoinRouter(route *gin.RouterGroup) {
	route.Any("/", func(context *gin.Context) {
		context.String(0, route.BasePath())
	})

	controller.CoinController(route.Group("coin"))
	controller.RatioController(route.Group("ratio"))

}
