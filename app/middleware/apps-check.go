package middleware

import (
	"github.com/kataras/iris"
	"im/app"
	"im/common"
	"im/service"
	"net/http"
	"strings"
)

//app token验证
func AppsCheck(ctx iris.Context){
	appToken:=strings.Trim(ctx.URLParam(app.GET_NAME_OF_APP_TOKEN)," ")
	if appToken ==""{
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(common.SendSad("not found "+app.GET_NAME_OF_APP_TOKEN+" in url"))
		return
	}

	appId,err := service.NewAppService().GetToken(appToken)
	if err!=nil{
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(common.SendSad("invalid "+app.GET_NAME_OF_APP_TOKEN))
		return
	}
	ctx.Values().Set(app.MIDDLEWARE_APP_ID_KEY,appId)
	ctx.Next()
}
