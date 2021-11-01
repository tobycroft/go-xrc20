package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/user/model/UserModel"
	"main.go/app/v1/wallet/model/UserAddressModel"
	"main.go/common/BaseController"
	"main.go/common/BaseModel/TokenModel"
	"main.go/tuuz"
	"main.go/tuuz/Base64"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
	"main.go/tuuz/Vali"
)

func AddressController(route *gin.RouterGroup) {
	route.Use(BaseController.CorsController())

	route.Any("create", address_without_create)
	route.Any("import", address_import)
	route.Any("memoric", address_memoric)
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

func address_without_create(c *gin.Context) {
	username, ok := Input.Post("wallet_name", c, true)
	if !ok {
		return
	}
	err := Vali.Length(username, 1, 12)
	if err != nil {
		RET.Fail(c, 400, err.Error(), "钱包长度不正确")
		return
	}
	password, ok := Input.Post("password", c, false)
	err = Vali.Length(password, 5, 19)
	if err != nil {
		RET.Fail(c, 400, err.Error(), "密码长度不正确")
		return
	}
	pass_notify, ok := Input.Post("pass_notify", c, true)
	if !ok {
		return
	}
	if len(UserModel.Api_find_byUsername(username)) > 0 {
		RET.Fail(c, 400, nil, "用户名已经被注册")
		return
	}
	invite_code, ok := Input.Post("invite_code", c, false)
	if !ok {
		return
	}
	if len(UserModel.Api_find_byUsername(username)) > 0 {
		RET.Fail(c, 400, nil, "用户名已经被注册")
		return
	}
	invite_user, err := Base64.Decode(invite_code)
	if err != nil {
		RET.Fail(c, 406, nil, "邀请码不正确")
		return
	}
	invite_data := UserModel.Api_find_byUsername(invite_user)
	if len(invite_data) < 1 {
		RET.Fail(c, 404, nil, "邀请人不存在")
		return
	}

	var usermodel UserModel.Interface
	db := tuuz.Db()
	db.Begin()
	usermodel.Db = db
	uid := usermodel.Api_insert(invite_data["id"], username, Calc.Md5(password), pass_notify, nil, "cn", Base64.Encode([]byte(username)))
	if uid != 0 {
		token := Calc.GenerateToken()
		if !TokenModel.Api_insert(uid, token, "app") {
			db.Rollback()
			RET.Fail(c, 401, nil, "token写入失败")
			return
		}
		words_sli, words := WordAction.App_gen_words()
		var useraddress UserAddressModel.Interface
		useraddress.Db = db
		if !useraddress.Api_insert("user", uid, "", words) {
			db.Rollback()
			RET.Fail(c, 500, nil, "地址插入失败")
			return
		}
		db.Commit()
		RET.Success(c, 0, nil, map[string]interface{}{
			"uid":       uid,
			"token":     token,
			"words_sli": words_sli,
			"words":     words,
		})
	} else {
		db.Rollback()
		RET.Fail(c, 500, nil, "用户创建失败")
	}
}

func address_import(c *gin.Context) {
	//恢复身份
	secret, ok := Input.Post("secret", c, true)
	if !ok {
		return
	}
	var aa AddressModel.Interface
	aa.Db = tuuz.Db()
	addr := aa.Api_find_secretCode(secret)
	//addr := aa.Api_find_secretCode(secret)
	//fmt.Println(AddressAction.App_address_secret_decrypt(addr["secret"].(string)))
	if len(addr) > 1 {
		sec, err := AddressAction.App_address_secret_decrypt(addr["secret"].(string))
		if err != nil {
			RET.Fail(c, 407, err.Error(), "私钥无法解密！")
			return
		}
		err = AddressAction.App_address_eccChecker(addr["address"].(string), sec, addr["ecc_code"].(string))
		if err != nil {
			RET.Fail(c, 407, err.Error(), "私钥无法通过风控检测！")
			return
		}
		var ur UserAddressModel.Interface
		ur.Db = tuuz.Db()
		ua := ur.Api_find_address(addr["address"].(string))
		if len(ua) < 1 {
			RET.Fail(c, 404, nil, "未查询到对应账户，建议使用助记词恢复")
			return
		}
		//todo:
		token := Calc.GenerateToken()
		if !TokenModel.Api_insert(ua["uid"], token, "app") {
			RET.Fail(c, 401, nil, "token写入失败")
			return
		}
		RET.Success(c, 0, map[string]interface{}{
			"uid":   ua["uid"],
			"token": token,
		}, nil)
	} else {
		RET.Fail(c, 404, nil, "非本系统私钥，无法导入")
	}
}

func address_memoric(c *gin.Context) {
	words, ok := Input.Post("words", c, false)
	if !ok {
		return
	}
	var ur UserAddressModel.Interface
	ur.Db = tuuz.Db()
	user := ur.Api_find_byWord(words)
	if len(user) > 0 {
		us := UserModel.Api_find(user["uid"])
		if len(us) < 1 {
			RET.Fail(c, 404, nil, "未找到用户")
			return
		}
		if us["active"].(int64) != 1 {
			RET.Fail(c, 401, nil, "用户无法登录")
			return
		}
		token := Calc.GenerateToken()
		if !TokenModel.Api_insert(user["uid"], token, "app") {
			RET.Fail(c, 401, nil, "token写入失败")
			return
		}
		RET.Success(c, 0, map[string]interface{}{
			"uid":   user["uid"],
			"token": token,
		}, nil)
	} else {
		RET.Fail(c, 404, nil, "未找到用户")
	}
}
