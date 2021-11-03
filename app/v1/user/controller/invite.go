package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InviteController(route *gin.RouterGroup) {
	route.Use(cors.Default())

	route.Any("get", info_share)
}
