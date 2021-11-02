package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/coin/action/RatioAction"
	"main.go/app/v1/coin/model/CoinRatioModel"
	"main.go/common/BaseController"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
)

func RatioController(route *gin.RouterGroup) {
	route.Use(BaseController.CorsController())

	route.Any("list", ratio_list)
	route.Any("get", ratio_get)
}

func ratio_list(c *gin.Context) {
	coins := RatioAction.App_Ratio_handler()
	RET.Success(c, 0, coins, nil)
}

func ratio_get(c *gin.Context) {
	rid, ok := Input.PostInt("rid", c)
	if !ok {
		return
	}
	coin := CoinRatioModel.Api_find(rid)
	if len(coin) > 1 {
		RET.Success(c, 0, coin, nil)
	} else {
		RET.Fail(c, 404, nil, nil)
	}
}
