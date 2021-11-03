package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/app/v1/user/action/UserInfoAction"
	"main.go/common/BaseModel/SystemParamModel"
	"main.go/tuuz/RET"
)

func InfoController(route *gin.RouterGroup) {
	route.Use(cors.Default())

	route.Any("get", info_get)
	route.Any("share", info_share)
}

func info_get(c *gin.Context) {
	uid := c.PostForm("uid")
	user := UserInfoAction.App_userinfo(uid)
	RET.Success(c, 0, user, nil)
}
func info_share(c *gin.Context) {
	share_url := SystemParamModel.Api_find_val("share_url")
	RET.Success(c, 0, share_url, nil)
}
