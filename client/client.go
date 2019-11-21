package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"im/client/common"
	"im/client/model"
	"im/client/service"
	"os"
	"strings"
)

func main() {
	var(
		a=&model.App{}
		reader = bufio.NewReader(os.Stdin)
		r = &model.Response{}
		//appAuthResponse []byte
		u = new(model.Users)
		tokenFile *os.File
	)

	RETRY:
	fmt.Fprintln(os.Stdout, "加载app授权token中... \n")

	//获取历史app token
	tokenFile,err:=os.OpenFile(common.TokenFileName,os.O_RDWR,0751)
	defer tokenFile.Close()
	if err!=nil{
		if os.IsPermission(err){
			panic("token文件打开失败,权限不足 "+err.Error())
		}
		if os.IsNotExist(err){
			tokenFile,err:=os.Create(common.TokenFileName)
			if err!=nil{
				panic("token存储文件创建失败 原因:"+err.Error())
			}
			err=tokenFile.Chmod(0751)
			tokenFile,_=os.OpenFile(common.TokenFileName,os.O_RDWR,0751)
		}
	}

	buf:=make([]byte,64)
	if l,_:=tokenFile.Read(buf);l==0{
		a.GetNewAppToken(tokenFile)
	}else{
		a.Token=buf
	}

	fmt.Fprintf(os.Stdout,"获取app token 成功! %s \n",a.Token)



	fmt.Fprintln(os.Stdout,"请输入聊天账号:")
	u.Account,_ = reader.ReadString('\n')
	fmt.Fprintln(os.Stdout,"请输入聊天密码:")
	u.Password,_ = reader.ReadString('\n')
	u.Account,u.Password = strings.TrimRight(u.Account,"\n"),strings.TrimRight(u.Password,"\n")
	u.Account,u.Password = strings.TrimRight(u.Account,"\r\n"),strings.TrimRight(u.Password,"\r\n")
	authToken:=common.Post(fmt.Sprintf("%s/users/token?X-Websocket-Header-X-APP-Token=%s",common.Host,a.Token),&u,common.Headers{"Content-Type":"application/json"})
	json.Unmarshal(authToken,r)
	if r.Code !=0{
		fmt.Fprintln(os.Stdout,"登录失败 原因:"+r.Msg)
		fmt.Fprintln(os.Stdout,"是否重新授权后重试: ( y/n )")
		retry,_:=reader.ReadString('\n');
		retry=strings.TrimRight(retry,"\n")
		retry=strings.TrimRight(retry,"\r\n")
		if retry=="y" || retry=="Y"{
			tokenFile.Truncate(0)
			goto RETRY
		}else{
			panic(r.Msg)
		}
	}
	u.Token=&r.Data
	fmt.Fprintln(os.Stdout,"登录成功,jwt:"+*u.Token)
	fmt.Fprintln(os.Stdout,"正在连接websocket")

	token,err:=jwt.Parse(*u.Token, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")),nil
	})
	if token.Claims == nil{
		fmt.Fprintln(os.Stdout,"JWT failure")
		return
	}
	user:=token.Claims.(jwt.MapClaims)

	service.NewChat().Connect(string(a.Token),*u.Token,uint(user["apps_id"].(float64)),uint64(user["id"].(float64)))
}