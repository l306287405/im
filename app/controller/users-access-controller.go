package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/app"
	"im/common"
	"im/service"
	"net/http"
)

type UsersAccessController struct {
	Session *sessions.Session
	Ctx iris.Context
}

type usersAccessPostParams struct {
	Account		string	`json:"account"`
	Password	string	`json:"password"`
	Nickname	string	`json:"nickname"`
}

//创建账号
func (c *UsersAccessController) Post(){
	var(
		user = &usersAccessPostParams{}
		err = c.Ctx.ReadJSON(user)
		appsId,_ = c.Ctx.Values().GetUint(app.MIDDLEWARE_APP_ID_KEY)
		result int64
	)

	if err!=nil{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendSad("参数缺失:"+err.Error()))
	}

	if user.Account=="" || user.Password=="" || user.Nickname==""{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendSad("参数缺失"))
		return
	}

	result,err =service.NewUserService().Create(appsId,user.Account,user.Password,user.Nickname)

	if err!=nil{
		c.Ctx.StatusCode(http.StatusInternalServerError)
		c.Ctx.JSON(common.SendCry("创建账号失败",err.Error()))
		return
	}
	c.Ctx.StatusCode(http.StatusCreated)
	c.Ctx.JSON(common.SendSmile(result,"创建账号成功"))
}