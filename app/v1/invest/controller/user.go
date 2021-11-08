package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/invest/model/InvestUserModel"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/RET"
)

func UserController(route *gin.RouterGroup) {

	route.Use(BaseController.LoginedController(), gin.Recovery())
	route.Any("get", user_get)
}

func user_get(c *gin.Context) {
	uid := c.PostForm("uid")
	var iu InvestUserModel.Interface
	iu.Db = tuuz.Db()
	investuser := iu.Api_find(uid)
	if len(investuser) > 0 {
		RET.Success(c, 0, 1, nil)
	} else {
		RET.Success(c, 0, 0, nil)
	}
}
