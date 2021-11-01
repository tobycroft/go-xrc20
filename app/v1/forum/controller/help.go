package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/app/v1/forum/model/ForumModel"
	"main.go/app/v1/forum/model/ForumThreadModel"
	"main.go/app/v1/index/model/HelpModel"
	"main.go/app/v1/index/model/SystemParamModel"
	"main.go/tuuz/Input"
	"main.go/tuuz/RET"
	"main.go/tuuz/Vali"
)

func HelpController(route *gin.RouterGroup) {
	route.Any("/", func(context *gin.Context) {
		context.String(0, route.BasePath())
	})

	route.Any("forum", help_forum)
	route.Any("cata", help_cata)
	route.Any("hot", help_hot)
	route.Any("search", help_search)
	route.Any("get", thread_get)
}

func help_forum(c *gin.Context) {
	fid, ok := Input.PostInt("fid", c)
	if !ok {
		return
	}
	limit, page, err := Input.PostLimitPage(c)
	if err != nil {
		return
	}
	//输出点击进入帮助板块的帖子
	datas := ForumThreadModel.Api_select(fid, limit, page)
	RET.Success(c, 0, datas, nil)
}

func help_search(c *gin.Context) {
	search_text, ok := Input.Post("search_text", c, true)
	if !ok {
		return
	}
	err := Vali.Length(search_text, 2, 16)
	if err != nil {
		RET.Fail(c, 400, err.Error(), nil)
		return
	}
	datas := HelpModel.Api_select()
	fids := []interface{}{}
	for _, data := range datas {
		fids = append(fids, data["fid"])
	}
	threads := ForumThreadModel.Api_like_in(search_text, fids)
	RET.Success(c, 0, threads, nil)
}

func help_hot(c *gin.Context) {
	limit, page, err := Input.PostLimitPage(c)
	if err != nil {
		return
	}
	help_hot_fid := SystemParamModel.Api_find_val("help_hot_fid")
	if help_hot_fid != nil {
		forum := ForumModel.Api_find(help_hot_fid)
		if len(forum) > 0 {
			datas := ForumThreadModel.Api_select(help_hot_fid, limit, page)
			RET.Success(c, 0, datas, nil)
		} else {
			RET.Fail(c, 404, nil, "未找到板块")
		}
	} else {
		RET.Fail(c, 400, nil, "未设定热门帮助板块")
	}
}

func help_cata(c *gin.Context) {
	datas := HelpModel.Api_select()
	fids := []interface{}{}
	for _, data := range datas {
		fids = append(fids, data["fid"])
	}
	forums := ForumModel.Api_select_in(fids)
	RET.Success(c, 0, forums, nil)
}
