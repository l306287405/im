package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"im/client/common"
	"os"
	"strings"
)

type App struct {
	Id string	`json:"key_id"`
	Secret string	`json:"key_secret"`
	Token []byte	`json:"token,omitempty"`
}

func (a *App) GetNewAppToken(tokenFile *os.File){
	var(
		reader = bufio.NewReader(os.Stdin)
		r = &Response{}
		appAuthResponse []byte
		request map[string]string
		headers=make(map[string]string)

		err error
	)

	fmt.Fprint(os.Stdout, "请输入app授权id... \n")

	a.Id,_ = reader.ReadString('\n')
	if a.Id==""{
		a.Id="wechat"
	}
	fmt.Fprint(os.Stdout, "请输入app授权secret... \n")
	a.Secret,_ = reader.ReadString('\n')
	if a.Secret==""{
		a.Secret="test"
	}
	fmt.Fprint(os.Stdout, "请求获取app调用token中... \n")
	a.Id,a.Secret=strings.TrimRight(a.Id,"\n"),strings.TrimRight(a.Secret,"\n")

	request= map[string]string{"key_id":a.Id,"key_secret":a.Secret}
	headers["Content-Type"]="application/json"

	appAuthResponse,err = common.HttpDo("GET",common.Host+"/apps/token",request,&headers)
	if err!=nil{
		fmt.Fprintln(os.Stdout,err.Error())
		return
	}
	fmt.Println(string(appAuthResponse))
	json.Unmarshal(appAuthResponse,r)

	if r.Code != 0{
		panic("app授权获取失败")
	}
	a.Token=[]byte(r.Data)
	_,err=tokenFile.Write(a.Token)
	if err!=nil{
		panic("token写入失败,原因:"+err.Error())
	}
}