package Middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"im/common"
	"im/model"
	"im/service"
	"net/http"
	"os"
	"strings"
)

//用户token校验中间件
func UsersVerify(c iris.Context){
	var(
	 	tokenStr string
	 	cachedToken string
	 	loginUser model.Users
	 	token *jwt.Token
		authHeaderParts []string
	 )

	tokenStr=c.GetHeader("Authorization")
	if tokenStr==""{
		c.StatusCode(http.StatusUnauthorized)
		c.JSON(common.SendCry("unauthorized"))
		return
	}
	authHeaderParts = strings.Split(tokenStr, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		c.StatusCode(http.StatusUnauthorized)
		c.JSON(common.SendCry("Authorization header format must be Bearer {token}"))
		return
	}

	token,err:=jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")),nil
	})

	if token.Claims == nil{
		c.StatusCode(http.StatusUnauthorized)
		c.JSON(common.SendCry("JWT failure"))
		return
	}
	user:=token.Claims.(jwt.MapClaims)
	cachedToken,err=service.NewUserService().GetCacheByUid(uint64(user["id"].(float64)))
	if err!=nil{
		c.StatusCode(http.StatusBadRequest)
		c.JSON(common.SendCry("JWT Matching failure"))
		return
	}
	if authHeaderParts[1] != cachedToken{
		c.StatusCode(http.StatusUnauthorized)
		c.JSON(common.SendCry("JWT Matching failure"))
		return
	}

	loginUser.Id=uint64(user["id"].(float64))
	loginUser.AppsId=uint(user["apps_id"].(float64))
	loginUser.Account=user["account"].(string)
	loginUser.Nickname=user["nickname"].(string)

	c.Values().Set("user",loginUser)
	c.Next()
}