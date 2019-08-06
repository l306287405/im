package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Host				  = "http://127.0.0.1:8080"
	Endpoint              = "ws://localhost:8080/websocket/echo"
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