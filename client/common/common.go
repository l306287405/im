package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Host				  = "http://10.0.0.12:8080"
	Endpoint              = "ws://10.0.0.12:8080/websocket/echo"
	DialAndConnectTimeout = 5 * time.Second
	TokenFileName		  = "token.txt"
	APP_NAME_OF_APP_TOKEN="X-Websocket-Header-X-APP-Token"

	//app校验通过后中间件传递给控制器的app_id所用key名
	MIDDLEWARE_APP_ID_KEY="APP_ID"
)

type Headers map[string]string

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data interface{}, h Headers) []byte {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err!=nil{
		fmt.Println(err.Error())
	}

	for k,v:=range h {
		req.Header.Add(k, v)
	}

	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return result
}

func HttpDo(method string,url string,jsonData *string,headers *map[string]string) ([]byte,error) {
	var(
		req = fasthttp.AcquireRequest()
		resp = fasthttp.AcquireResponse()
		result = []byte("")
		err error
	)

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(method)
	if headers!=nil{
		for k,v:=range *headers{
			req.Header.Set(k,v)
		}
	}

	req.SetRequestURI(url)

	if jsonData!=nil{
		req.SetBody([]byte(*jsonData))
	}


	if err = fasthttp.Do(req, resp); err != nil {
		return result,err
	}

	result = resp.Body()
	return result,nil
}