package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"im/app"
	"im/client/common"
	model2 "im/client/model"
	"im/client/service"
	"im/model"
	"os"
	"strings"
)

func main() {
	var(
		a=&model2.App{}
		reader = bufio.NewReader(os.Stdin)
		r = &model2.Response{}
		//appAuthResponse []byte
		u = &model.Users{}
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
	authToken:=common.Post(common.Host+"/users",&u,common.Headers{"Content-Type":"application/json",app.HEADER_NAME_OF_APP_TOKEN:string(a.Token)})
	json.Unmarshal(authToken,r)
	if r.Code !=0{
		fmt.Fprintln(os.Stdout,"登录失败 原因:"+r.Msg)
		fmt.Fprintln(os.Stdout,"是否重新授权后重试: ( y/n )")
		retry,_:=reader.ReadString('\n');
		retry=strings.TrimRight(retry,"\n")
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

	service.NewChat().Connect("Bearer "+*u.Token)
}