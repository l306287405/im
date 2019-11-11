package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"sync"
)

var (
	instance *oss.Client
	once sync.Once
)

func GetOSS() *oss.Client {
	once.Do(func() {
		// Endpoint以杭州为例，其它Region请按实际情况填写。
		endpoint := "http://"+os.Getenv("OSS_ENDPOINT")
		// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
		accessKeyId := os.Getenv("ALI_ACCESS_ID")
		accessKeySecret := os.Getenv("ALI_ACCESS_SECRET")

		client, err := oss.New(endpoint,accessKeyId,accessKeySecret)
		if err!=nil{
			panic("oss对象连接失败,"+err.Error())
		}

		instance = client
	})
	return instance
}

func GetBucket() *oss.Bucket {
	client:=GetOSS()
	bucket,err:=client.Bucket(os.Getenv("OSS_BUCKET"))
	if err!=nil{
		panic("所选bucket不存在:"+err.Error())
	}
	return bucket
}
