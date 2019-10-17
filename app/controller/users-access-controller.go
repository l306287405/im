package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/app"
	"im/common"
	"im/model"
	"im/service"
	"net/http"
	"strings"
)

type UsersAccessController struct {
	Session *sessions.Session
	Ctx iris.Context
}

//登录
func (c *UsersAccessController) Post(){
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
func (c *UsersAccessController) Delete(){
	var(
		token string
	)

	jwtStr := c.Ctx.Values().Get("jwt")
	if jwtStr == nil{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendCry("用户身份校验失败"))
		return
	}
	user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)
	cachedToken,err:=service.NewUserService().GetCacheByUid(uint64(user["id"].(float64)))
	if err!=nil{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("JWT Matching failure"))
		return
	}
	if token != cachedToken{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendCry("JWT Matching failure"))
		return
	}

	err=service.NewUserService().DelCacheOfToken(uint64(user["id"].(float64)))

	if err!=nil{
		c.Ctx.StatusCode(http.StatusInternalServerError)
		c.Ctx.JSON(common.SendCry("退出登录失败",err.Error()))
		return
	}
	c.Ctx.JSON(common.SendSmile(nil,"退出登录成功"))
}

//创建账号
func (c *UsersAccessController) Put(){
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