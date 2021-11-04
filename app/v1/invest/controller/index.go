package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/RET"
)

func IndexController(route *gin.RouterGroup) {

	route.Use(BaseController.CorsController())
	//route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("test", index_test)
	route.Any("delete", index_delete)
	route.Any("delete_all", index_delete_all)
}

func index_test(c *gin.Context) {
	db := tuuz.Db()
	db.Execute("UPDATE `fbcct`.`invest_order` SET `last_exec` = 0")
	//Invest.Invest_function()
	RET.Success(c, 0, nil, "执行需要时间，等个几十秒吧")
}

func index_delete(c *gin.Context) {
	tuuz.Db().Table("invest_order").Truncate()
	tuuz.Db().Table("invest_user").Truncate()
}

func index_delete_all(c *gin.Context) {
	tuuz.Db().Table("invest_order").Truncate()
	tuuz.Db().Table("invest_user").Truncate()
	tuuz.Db().Table("fb_balance").Truncate()
	tuuz.Db().Table("fb_balance_record").Truncate()
}
