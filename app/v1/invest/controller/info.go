package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/invest/model/InvestOrderModel"
	"main.go/app/v1/invest/model/InvestUserModel"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/Calc"
	"main.go/tuuz/RET"
)

func InfoController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("get", info_get)
	route.Any("create", info_create)
}

func info_get(c *gin.Context) {
	uid := c.PostForm("uid")
	var iu InvestUserModel.Interface
	iu.Db = tuuz.Db()
	investuser := iu.Api_find(uid)
	var iv InvestOrderModel.Interface
	iv.Db = tuuz.Db()
	amount, _, _, _ := iv.Api_sum_byUid(uid)
	if len(investuser) > 0 {
		data := map[string]interface{}{
			//锁仓:就是分红-锁仓额度
			"lock": Calc.ToDecimal(investuser["freeze_amount"]),
			//锁仓额度
			"lock_amount": Calc.ToDecimal(investuser["lock_amount"]),
			//待释放额度
			"on_release": Calc.ToDecimal(investuser["on_release"]),
			//分红权
			"amount":       Calc.ToDecimal(investuser["amount"]),
			"level_amount": Calc.ToDecimal(investuser["level_amount"]).Sub(Calc.ToDecimal(amount)),
			"level":        Calc.ToDecimal(investuser["level"]),
		}
		RET.Success(c, 0, data, nil)
	} else {
		RET.Fail(c, 404, nil, "请先访问create接口创建一个账户")
	}
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
