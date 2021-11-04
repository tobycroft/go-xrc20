package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/invest/model/InvestOrderModel"
	"main.go/app/v1/user/action/UserInfoAction"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
)

func RecordController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("fenhong", record_fenhong)

}

func record_fenhong(c *gin.Context) {
	uid := c.PostForm("uid")
	limit, page, err := Input.PostLimitPage(c)
	if err != nil {
		return
	}
	var inv InvestOrderModel.Interface
	inv.Db = tuuz.Db()
	order := inv.Api_find_first(uid)
	if len(order) < 1 {
		RET.Fail(c, 0, nil, nil)
		return
	}
	datas := inv.Api_select_group_orderIdUid_byPidandDate(uid, order["date"], limit, page)
	for i, data := range datas {
		data["user_info"] = UserInfoAction.App_userinfo(data["uid"])
		datas[i] = data
	}
	RET.Success(c, 0, datas, nil)
}
