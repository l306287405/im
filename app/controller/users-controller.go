package controller

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/dao"
	"im/model"
	"net/http"
	"strings"
)

type UsersController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c *UsersController) Get(){
	var(
		userToken = model.Users{}
		user=&model.Users{}
	)
	userToken=c.Ctx.Values().Get("user").(model.Users)

	user=dao.NewUsersDao().Info(userToken.Id)
	c.Ctx.JSON(common.SendSmile(user))
	return

}

type putParams struct {
	Nickname *string	`json:"nickname"`
	Status *byte	`json:"status"`
}

func (c *UsersController) Put(){
	var(
		user = model.Users{}
		params = &putParams{}
		changedStr []string
		err=c.Ctx.ReadJSON(params)
	)
	if err!=nil{
		goto PARAMS_ERR
	}
	user=c.Ctx.Values().Get("user").(model.Users)

	if nickname:=params.Nickname;nickname!=nil{
		user.Nickname=strings.TrimSpace(*nickname)
		if len(user.Nickname)<3{
			println(user.Nickname)
			err=errors.New("昵称太短,长度请大于2.")
			goto PARAMS_ERR
		}
		changedStr=append(changedStr,"nickname")
	}
	if status:=params.Status;status!=nil{
		user.Status=*status
		if user.Status!=0 && user.Status!=1 {
			err=errors.New("用户状态不在可取值范围内.")
			goto PARAMS_ERR
		}
		changedStr=append(changedStr,"status")
	}

	err=dao.NewUsersDao().UpdateById(user.Id,&user,changedStr...)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(1,"更新成功"))

	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return
}