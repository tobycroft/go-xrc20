package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/common/BaseController"
)

func PaymentController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

}

func payment_buy(c *gin.Context) {

}
