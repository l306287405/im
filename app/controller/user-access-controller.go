package controller

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/app"
	"im/common"
	"im/model"
	"im/service"
	"im/service/orm"
	"net/http"
	"strings"
)

type UserAccessController struct {
	Session *sessions.Session
	Ctx iris.Context
}

//登录
func (c *UserAccessController) Post(){
	var(
		appId,_ = c.Ctx.Values().GetUint(app.MIDDLEWARE_APP_ID_KEY)
		u = &model.Users{}
	)
	c.Ctx.ReadJSON(u)
	if u.Account == "" || u.Password == "" {
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("参数缺失"))
		return
	}

	token,err := service.NewUserService().Login(appId,u.Account,[]byte(u.Password))
	if err!=nil{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("登录失败:"+err.Error()))
		return
	}
	c.Ctx.JSON(common.SendSmile(token))
	return
}

//退出登录
func (c *UserAccessController) Delete(){
	var(
		token string
		db *xorm.Engine
		users model.Users
	)

	token=c.Ctx.GetHeader("token")

	if token==""{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("参数缺失"))
		return
	}

	users.Token=&token
	db=orm.GetDB()
	r,e:=db.Cols("token").Where("token=?", token).Update(users)
	if e!=nil{
		c.Ctx.StatusCode(http.StatusInternalServerError)
		c.Ctx.JSON(common.SendCry("退出登录失败",e.Error()))
		return
	}
	e=service.NewUserService().DelCacheOfToken(token)
	if e!=nil{
		c.Ctx.StatusCode(http.StatusInternalServerError)
		c.Ctx.JSON(common.SendCry("退出登录失败",e.Error()))
	}
	c.Ctx.JSON(common.SendSmile(r,"退出登录成功"))
}

//创建账号
func (c *UserAccessController) Put(){
	var(
		account = strings.Trim(c.Ctx.PostValue("account")," ")
		password = strings.Trim(c.Ctx.PostValue("password")," ")
		nickname = strings.Trim(c.Ctx.PostValue("nickname")," ")
		appsId,err = c.Ctx.Values().GetUint(app.MIDDLEWARE_APP_ID_KEY)
		result int64
	)

	if account==" " || password==" " || nickname==" "{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendSad("参数缺失"))
		return
	}

	result,err =service.NewUserService().Create(appsId,account,password,nickname)

	if err!=nil{
		c.Ctx.StatusCode(http.StatusInternalServerError)
		c.Ctx.JSON(common.SendCry("创建账号失败",err.Error()))
		return
	}
	c.Ctx.StatusCode(http.StatusCreated)
	c.Ctx.JSON(common.SendSmile(result,"创建账号成功"))
}