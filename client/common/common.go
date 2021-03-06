package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	Host				  = "http://10.0.0.12:8080"
	Endpoint              = "ws://10.0.0.12:8080/websocket/echo"
	DialAndConnectTimeout = 5 * time.Second
	TokenFileName		  = "token.txt"
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

func HttpDo(method string,urlStr string,jsonData interface{},headers *map[string]string) ([]byte,error) {
	var(
		req = fasthttp.AcquireRequest()
		resp = fasthttp.AcquireResponse()
		result = []byte("")
		requestStrs []string
		requestStr string
		requestByte []byte

		interval string

		err error
	)

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(strings.ToUpper(method))
	if headers!=nil{
		for k,v:=range *headers{
			req.Header.Set(k,v)
		}
	}

	if jsonData!=nil{
		if method=="GET"{
			r:=reflect.ValueOf(jsonData).MapRange()
			for{
				if !r.Next() {
					break
				}
				requestStrs = append(requestStrs,r.Key().String()+"="+fmt.Sprint(r.Value().Interface()))
			}
			requestStr=strings.Join(requestStrs,"&")

			if strings.Contains(urlStr,"?"){
				interval="&"
			}else{
				interval="?"
			}
			urlStr+=interval+requestStr
			fmt.Println(urlStr)
		}else{
			requestByte,_=json.Marshal(jsonData)
			req.SetBody(requestByte)
		}
	}

	req.SetRequestURI(urlStr)

	if err = fasthttp.Do(req, resp); err != nil {
		return result,err
	}

	result = resp.Body()
	return result,nil
}