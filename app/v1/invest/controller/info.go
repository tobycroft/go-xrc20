package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/invest/model/InvestUserModel"
	"main.go/common/BaseController"
	"main.go/common/BaseModel/SystemParamModel"
	"main.go/tuuz"
	"main.go/tuuz/RET"
)

func InfoController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("get", info_get)
	route.Any("create", info_create)
}

func info_get(c *gin.Context) {
	invest_total_num := SystemParamModel.Api_find_val("invest_total_num")

	RET.Success(c, 0, map[string]interface{}{
		"total": invest_total_num,
		"has":   150,
	}, nil)
}

func info_create(c *gin.Context) {
	uid := c.PostForm("uid")
	var iv InvestUserModel.Interface
	iv.Db = tuuz.Db()
	if len(iv.Api_find(uid)) < 1 {
		if iv.Api_insert(uid) {
			RET.Success(c, 0, nil, "创建成功")
		} else {
			RET.Fail(c, 500, nil, nil)
		}
	} else {
		RET.Success(c, 0, nil, "")
	}
}
