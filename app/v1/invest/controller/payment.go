package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/common/BaseController"
	"main.go/tuuz/Input"
)

func PaymentController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

}

func payment_buy(c *gin.Context) {
	uid := c.PostForm("uid")
	from, ok := Input.Post("from", c, false)
	if !ok {
		return
	}
	to, ok := Input.Post("to", c, false)
	if !ok {
		return
	}
	amount, ok := Input.PostDecimal("amount", c)
	if !ok {
		return
	}
	mode, ok := Input.PostInt64("mode", c)
	if !ok {
		return
	}

}
