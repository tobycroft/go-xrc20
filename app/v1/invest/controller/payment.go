package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"main.go/app/cron/Invest"
	"main.go/app/v1/balance/action/BalanceAction"
	"main.go/app/v1/coin/action/CoinAction"
	"main.go/common/BaseModel/SystemParamModel"

	"main.go/app/v1/invest/action/PaymentAction"
	"main.go/app/v1/invest/model/InvestOrderModel"
	"main.go/app/v1/invest/model/InvestUserModel"
	"main.go/app/v1/user/model/UserModel"
	"main.go/common/BaseController"
	"main.go/tuuz"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
	"time"
)

func PaymentController(route *gin.RouterGroup) {
	route.Use(BaseController.CorsController())
	route.Use(BaseController.LoginedController(), gin.Recovery())

	route.Any("bill", payment_bill)
	route.Any("buy", payment_buy)
}

func payment_bill(c *gin.Context) {
	kvs := SystemParamModel.Api_KV()
	buy_credit := kvs["buy_credit"]
	buy_coin := kvs["buy_coin"]
	buy_reward := kvs["buy_reward"]
	buy_reward_cid := kvs["buy_reward_cid"]
	coins := CoinAction.App_coin()
	RET.Success(c, 0, map[string]interface{}{
		"price":            buy_credit,
		"coin":             buy_coin,
		"coin_info":        coins[Calc.Any2Int64(buy_coin)],
		"num":              1,
		"reward":           buy_reward,
		"reward_coin_info": coins[Calc.Any2Int64(buy_reward_cid)],
	}, nil)
}

func payment_buy(c *gin.Context) {
	uid := c.PostForm("uid")
	password, ok := Input.Post("password", c, true)
	if !ok {
		return
	}
	//购买多单
	numbers, ok := Input.PostDecimal("numbers", c)
	if !ok {
		return
	}
	numbers = numbers.Round(0)
	if numbers.LessThanOrEqual(decimal.Zero) {
		RET.Fail(c, 400, nil, "购买数量应该大于0")
	}
	user := UserModel.Api_find(uid)
	if Calc.Md5(password) != user["password"] {
		RET.Fail(c, 401, nil, "密码错误")
	} else {
		//Todo:开始支付流程
		order_id := Calc.GenerateOrderId()
		kvs := SystemParamModel.Api_KV()
		buy_coin, berr := Calc.Any2Int64_2(kvs["buy_coin"])
		buy_credit := kvs["buy_credit"]
		buy_reward_cid := kvs["buy_reward_cid"]
		buy_reward := kvs["buy_reward"]

		coins := CoinAction.App_coin()

		if buy_credit == "" || buy_reward_cid == "" || berr != nil || coins[buy_coin] == nil {
			RET.Fail(c, 406, nil, "buy_credit设定不正确")
			return
		}

		rebuy_lock_amount, err1 := Calc.Any2Float64_2(kvs["rebuy_lock_amount"])
		pid_reward_freeze_amount, err2 := Calc.Any2Float64_2(kvs["pid_reward_freeze_amount"])
		pid_reward_on_release, err3 := Calc.Any2Float64_2(kvs["pid_reward_on_release"])
		if err1 != nil || err2 != nil || err3 != nil {
			RET.Fail(c, 304, nil, "系统参数设定错误")
			return
		}
		buy_amount := Calc.ToDecimal(buy_credit).Mul(numbers).Abs()
		db := tuuz.Db()
		db.Begin()
		var bal BalanceAction.Interface
		bal.Db = db
		err := bal.App_single_balance(uid, buy_coin, 30, order_id, buy_amount.Abs().Neg(), "", "", "")
		if err != nil {
			db.Rollback()
			RET.Fail(c, 300, nil, err.Error())
			return
		}

		//a_num := invest.Api_count_byUid(user["pid"])
		//a_lock_amount_plus := PaymentAction.App_calc_amount(a_num)
		//如果B的锁仓不为0，就调用用lock-b_lock_amount_plus，大于0就写入锁仓，小于0就写入lockamount
		var investuser InvestUserModel.Interface
		investuser.Db = db
		var investorder InvestOrderModel.Interface
		investorder.Db = db
		var pay PaymentAction.Interface
		pay.Db = db

		//先把人准备好
		_, err = pay.App_findInvestUser(user["id"])
		if err != nil {
			db.Rollback()
			RET.Fail(c, 500, nil, err.Error())
			return
		}
		_, err = pay.App_findInvestUser(user["pid"])
		if err != nil {
			db.Rollback()
			RET.Fail(c, 500, nil, err.Error())
			return
		}

		//填写这个数据的时候先查一下我的权限
		for i := 0; i < int(numbers.IntPart()); i++ {
			num_pid := investorder.Api_count_byUid(user["pid"])
			am := float64(0)
			if num_pid > 0 {
				am = 4000
			}

			if !investorder.Api_insert(user["id"], user["pid"], order_id, buy_credit, 4000, am, 0, time.Now().Unix()+7*86400) {
				db.Rollback()
				RET.Fail(c, 500, nil, "无法插入invest数据条")
				return
			}

			/*
				b的变动结束后，需要把a的数据整个提出来查看影响的部分
				b的变动会影响到a的分红权和锁仓，分红权直接x用户倍数，锁仓加上去
			*/

			//查看之前已经参加过几次了（算分红权数量）
			//B自己购买了几次，影响给b加锁仓额度lock_amount

			/*
				1.给自己B这个uid，lock_amount加锁仓额度3000（4000），并且给fbcc字段+2000
				2.大多数情况下，锁仓额度会被锁仓对抵掉导致变为0
				3.用B的pid找到A的uid，给A的lock锁仓加3000（4000），给A的amount加4000（固定）
				4.如果B自己的lock不为0，lock_amount字段不为0，则用lock-lock_amount，值填给on_release
				5.如果A自己的lock（肯定不为0）不为0，lock_amount字段还有，则lock-lock_amount，值给on_release
			*/

			//先处理a的信息，把该给的积分给a
			//算下要给a加多少锁仓4000，如果a有320u，就下一步会结算3000走进待释放
			//新版改成释放+3000，锁仓+1000

			if num_pid > 0 {
				if !investuser.Api_incr_freezeAmount(user["pid"], pid_reward_freeze_amount) {
					db.Rollback()
					RET.Fail(c, 500, nil, "Api_incr_freezeAmount错误")
					return
				}
				if !investuser.Api_incr_onRelease(user["pid"], pid_reward_on_release) {
					db.Rollback()
					RET.Fail(c, 500, nil, "Api_incr_onRelease错误")
					return
				}
				if !investuser.Api_incr_amount(user["pid"], 1) {
					db.Rollback()
					RET.Fail(c, 500, nil, "Api_incr_amount错误")
					return
				}

			}
			//开始处理b的信息

			//先查下这是b第几次购买320u

			num := investorder.Api_count_byUid(user["id"])
			//算出这把能解锁多少
			//unlock_amount := PaymentAction.App_calc_amount(num - 1)
			if num-1 > 0 {
				//b因为购买了理财，只会影响锁仓额度部分,第一次加3000第二次加4000
				//新b因为购买了理财，只会影响锁仓额度部分,第一次加0第二次加4000
				if !investuser.Api_incr_lockAmount(user["id"], rebuy_lock_amount) {
					db.Rollback()
					RET.Fail(c, 500, nil, "Api_incr_lockAmount错误")
					return
				}
			}
		}
		//将a下的锁仓释放进行平衡，会自动调整onrelease,计算分红总量
		err = pay.App_paybalance(user["pid"])
		if err != nil {
			db.Rollback()
			RET.Fail(c, 300, nil, err.Error())
			return
		}
		//b平掉自己的仓，统计释放数量,计算分红总量
		err = pay.App_paybalance(user["id"])
		if err != nil {
			db.Rollback()
			RET.Fail(c, 300, nil, err.Error())
			return
		}
		//赠送fbcct2000枚
		err = bal.App_single_balance(uid, buy_reward_cid, 101, order_id, Calc.ToDecimal(buy_reward).Mul(numbers), "增加认购权", "USDT购买的FBCC奖励,购买次数:"+numbers.String(), "")
		if err != nil {
			db.Rollback()
			RET.Fail(c, 300, nil, err.Error())
			return
		}

		var dynamic Invest.Interface
		dynamic.Db = db
		dynamic.Invest_dynamic_usdt(coins, buy_amount, uid, order_id)

		db.Commit()
		RET.Success(c, 0, nil, nil)

		level_refresh_no_delay := kvs["level_refresh_no_delay"]
		if level_refresh_no_delay == "1" {
			go func() {
				Invest.Level_data()
			}()
		}
	}
}
