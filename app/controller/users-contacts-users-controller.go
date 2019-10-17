package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
)

type UsersContactsUsersController struct {
	Session *sessions.Session
	Ctx iris.Context
}


//登录
func (c *UsersContactsUsersController) Post(){

	c.Ctx.JSON(common.SendSmile("hahaha"))
	return
}