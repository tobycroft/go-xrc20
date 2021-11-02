package RatioAction

import "main.go/app/v1/coin/model/CoinRatioModel"

func App_Ratio_handler() map[int64]interface{} {
	coins := CoinRatioModel.Api_select()
	arr := map[int64]interface{}{}
	for _, coin := range coins {
		arr[coin["id"].(int64)] = coin
	}
	return arr
}
