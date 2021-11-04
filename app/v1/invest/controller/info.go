package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/common/BaseController"
	"main.go/common/BaseModel/SystemParamModel"
	"main.go/tuuz/RET"
)

func InfoController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("get", info_get)
}

func info_get(c *gin.Context) {
	invest_total_num := SystemParamModel.Api_find_val("invest_total_num")

	RET.Success(c, 0, map[string]interface{}{
		"total": invest_total_num,
		"has":   150,
	}, nil)
}
