package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/service"
	"strings"
)

type AppsController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c AppsController) Post(){
	var(
		keyId = strings.Trim(c.Ctx.PostValue("key_id")," ")
		keySecret = strings.Trim(c.Ctx.PostValue("key_secret")," ")
	)

	if keyId=="" || keySecret==""{
		c.Ctx.JSON(common.SendCry("参数缺失"))
		return
	}

	token, err := service.NewAppService().Token(keyId, keySecret)
	if err!=nil{
		c.Ctx.JSON(common.SendSad("获取应用token失败 原因:"+err.Error()))
		return
	}
	c.Ctx.JSON(common.SendSmile(token))
	return

}