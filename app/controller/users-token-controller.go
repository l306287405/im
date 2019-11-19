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
	"os"
	"strings"
)

type UsersTokenController struct {
	Session *sessions.Session
	Ctx iris.Context
}

//登录
func (c *UsersTokenController) Post(){
	var(
		appId uint
		appToken string
		u = &model.Users{}
		token *string
		err error
	)

	appToken=strings.Trim(c.Ctx.URLParam(app.GET_NAME_OF_APP_TOKEN)," ")
	if appToken ==""{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendSad("not found "+app.GET_NAME_OF_APP_TOKEN+" in url"))
		return
	}

	appId,err = service.NewAppService().GetToken(appToken)
	if err!=nil{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendSad("invalid "+app.GET_NAME_OF_APP_TOKEN))
		return
	}

	c.Ctx.ReadJSON(u)
	if u.Account == "" || u.Password == "" {
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("参数缺失"))
		return
	}

	token,err = service.NewUserService().Login(appId,u.Account,[]byte(u.Password))
	if err!=nil{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("登录失败:"+err.Error()))
		return
	}
	c.Ctx.JSON(common.SendSmile(token))
	return
}

//退出登录
func (c *UsersTokenController) Delete(){
	var(
		token *jwt.Token
		jwtStr string
		user jwt.MapClaims
		err error
	)

	jwtStr=c.Ctx.URLParam(app.GET_NAME_OF_JWT_TOKEN)
	if jwtStr==""{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendCry("unauthorized of jwt"))
		return
	}
	token,err=jwt.Parse(jwtStr, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")),nil
	})

	if err!=nil{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendCry("JWT failure"))
		return
	}

	if token.Claims == nil{
		c.Ctx.StatusCode(http.StatusUnauthorized)
		c.Ctx.JSON(common.SendCry("JWT failure"))
		return
	}


	user = token.Claims.(jwt.MapClaims)
	cachedToken,err:=service.NewUserService().GetCacheByUid(uint64(user["id"].(float64)))
	if err!=nil{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendCry("JWT Matching failure"))
		return
	}
	if jwtStr != cachedToken{
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
