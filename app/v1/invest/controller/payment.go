package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/coin/model/CoinModel"
	"main.go/app/v1/invest/model/InvestModeModel"
	"main.go/app/v1/invest/model/InvestOrderModel"
	"main.go/app/v1/user/model/UserModel"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
)

func PaymentController(route *gin.RouterGroup) {
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("buy", payment_buy)
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
	Type, ok := Input.PostIn("type", c, []string{"eth", "trc"})
	if !ok {
		return
	}
	contract, ok := Input.Post("contract", c, false)
	if !ok {
		return
	}
	coin := CoinModel.Api_find_byTypeAndContract(Type, contract)
	if len(coin) < 1 {
		RET.Fail(c, 404, nil, "未找到coin")
		return
	}
	user := UserModel.Api_find(uid)
	if len(user) < 1 {
		RET.Fail(c, 403, nil, "未找到用户")
		return
	}
	investmode := InvestModeModel.Api_find(mode)
	if len(investmode) < 1 {
		RET.Fail(c, 404, nil, "未找到投资模式")
		return
	}
	db := tuuz.Db()
	var iv InvestOrderModel.Interface
	iv.Db = db
	data := iv.Api_find_compelete(uid)
	if len(data) > 0 {
		RET.Fail(c, 407, nil, "前一单未完成，请等待前一单完成或失败")
		return
	}
	order_id := Calc.GenerateOrderId()
	if !iv.Api_insert(uid, user["pid"], coin["id"], mode, order_id, amount, from, to, "", 0) {
		RET.Fail(c, 500, nil, nil)
		return
	}

}