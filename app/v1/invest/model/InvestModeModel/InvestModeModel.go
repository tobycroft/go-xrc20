package InvestModeModel

import (
	"github.com/tobycroft/gorose-pro"
	"main.go/tuuz"
	"main.go/tuuz/Log"
)

const table = "invest_mode"

func Api_find(generation interface{}) gorose.Data {
	db := tuuz.Db().Table(table)
	db.Where("generation", "=", generation)
	ret, err := db.Find()
	if err != nil {
		Log.Dbrr(err, tuuz.FUNCTION_ALL())
		return nil
	} else {
		return ret
	}
}

func Api_val(generation interface{}) interface{} {
	db := tuuz.Db().Table(table)
	db.Where("generation", "=", generation)
	ret, err := db.Value("amount")
	if err != nil {
		Log.Dbrr(err, tuuz.FUNCTION_ALL())
		return nil
	} else {
		return ret
	}
}
