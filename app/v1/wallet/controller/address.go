package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/user/model/UserModel"
	"main.go/app/v1/wallet/model/UserAddressModel"
	"main.go/common/BaseController"
	"main.go/common/BaseModel/TokenModel"
	"main.go/tuuz"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
	"main.go/tuuz/Vali"
)

func AddressController(route *gin.RouterGroup) {
	route.Use(BaseController.CorsController())

	route.Any("create", address_create)
}

func address_create(c *gin.Context) {
	address, ok := Input.Post("address", c, true)
	if !ok {
		return
	}
	password, ok := Input.Post("password", c, false)
	err := Vali.Length(password, 3, 19)
	if err != nil {
		RET.Fail(c, 400, err.Error(), "密码长度不符合要求")
		return
	}
	invite_code, ok := Input.Post("invite_code", c, false)
	if !ok {
		return
	}

	var ua UserAddressModel.Interface
	ua.Db = tuuz.Db()
	ua.Api_find_address(address)
	invite_data := UserModel.Api_find_byUsername(invite_code)
	if len(invite_data) < 1 {
		RET.Fail(c, 404, nil, "邀请人不存在")
		return
	}
	var usermodel UserModel.Interface
	db := tuuz.Db()
	db.Begin()
	usermodel.Db = db
	uid := usermodel.Api_insert(invite_data["id"], address, Calc.Md5(password), "", nil, "cn", address)
	if uid != 0 {
		token := Calc.GenerateToken()
		if !TokenModel.Api_insert(uid, token, "app") {
			db.Rollback()
			RET.Fail(c, 401, nil, "token写入失败")
			return
		}
		var useraddress UserAddressModel.Interface
		useraddress.Db = db
		if !useraddress.Api_insert("eth", uid, address, "") {
			db.Rollback()
			RET.Fail(c, 500, nil, "地址插入失败")
			return
		}
		db.Commit()
		RET.Success(c, 0, nil, map[string]interface{}{
			"uid":     uid,
			"token":   token,
			"address": address,
		})
	} else {
		db.Rollback()
		RET.Fail(c, 500, nil, "用户创建失败")
	}
}
