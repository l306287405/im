package Middleware

import (
	"github.com/kataras/iris"
	"im/common"
	"im/service"
	"net/http"
	"strings"
)

//app token验证
func AppsCheck(ctx iris.Context){
	appToken:=strings.Trim(ctx.GetHeader("APP-TOKEN")," ")
	if appToken ==""{
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(common.SendSad("not found APP-TOKEN in header"))
		return
	}

	appToken,err := service.NewAppService().GetCache(appToken)
	if err!=nil{
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(common.SendSad("invalid APP-TOKEN"))
		return
	}
	ctx.Next()
}
