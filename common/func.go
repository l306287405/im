package common

import (
	"strconv"
	"strings"
	"time"
)

const(
	SUCCESS byte=0
	WARNING byte=1
	ERROR	byte=2
)

type Response *map[string]interface{}

func DateTime(t time.Time)string{
	return t.Format("2006-01-02 15:04:05")
}

func SendSmile(data interface{},args ...interface{}) Response{
	result := make(map[string]interface{})
	for _,arg:=range args{
		switch arg.(type) {
		case string:
			result["msg"]=arg
		case byte:
			result["code"]=arg
		}
	}

	result["data"]=data
	if _,ok := result["msg"];!ok{
		result["msg"]="获取数据成功"
	}
	if _,ok := result["code"];!ok{
		result["code"]=SUCCESS
	}

	return &result
}

func SendSad(msg string,args ...interface{}) Response{
	result := make(map[string]interface{})
	for _,arg:=range args{
		_,ok := arg.(byte)
		if ok{
			result["code"]=arg
		}else {
			result["data"]=arg
		}
	}
	result["msg"]=msg
	if _,ok := result["data"];!ok{
		result["data"]=nil
	}
	if _,ok := result["code"];!ok{
		result["code"]=WARNING
	}
	return &result
}


func SendCry(msg string,args ...interface{}) Response{
	result := make(map[string]interface{})
	for _,arg:=range args{
		_,ok := arg.(byte)
		if ok{
			result["code"]=arg
		}else {
			result["data"]=arg
		}
	}
	result["msg"]=msg
	if _,ok := result["data"];!ok{
		result["data"]=nil
	}
	if _,ok := result["code"];!ok{
		result["code"]=ERROR
	}
	return &result
}

//解析唯一账号至appid与账号
func ParseAccount(account string) (uint,string){
	s:=strings.Split(account,"_")
	if len(s)!=2{
		return 0,""
	}
	appid,err:=strconv.Atoi(s[0])
	if err!=nil{
		return 0,""
	}

	return uint(appid),s[1]
}