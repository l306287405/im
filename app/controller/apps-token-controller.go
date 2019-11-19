package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/model"
	"im/service"
)

type AppsTokenController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c AppsTokenController) Get(){
	var(
		app=&model.Apps{}
		err error
	)
	err=c.Ctx.ReadJSON(app)
	if err!=nil{
		c.Ctx.JSON(common.SendCry("请求参数获取失败 原因:"+err.Error()))
		return
	}

	if app.KeyId=="" || app.KeySecret==""{
		c.Ctx.JSON(common.SendCry("参数缺失"))
		return
	}

	token, err := service.NewAppService().Token(app.KeyId, app.KeySecret)
	if err!=nil{
		c.Ctx.JSON(common.SendSad("获取应用token失败 原因:"+err.Error()))
		return
	}
	c.Ctx.JSON(common.SendSmile(token))
	return

}