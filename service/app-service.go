package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/cache"
	"im/service/orm"
	"log"
	"time"
)

type AppService struct {
	apps model.Apps
	db *xorm.Engine
	redis *redis.Client
}

func NewAppService() *AppService{
	return &AppService{db:orm.GetDB(),redis:cache.Get()}
}

//获取app 的访问token
func (s *AppService) Token(keyId string,keySecret string) (token string,err error){

	has, err := s.db.Cols("id").Where("key_id=?", keyId).And("key_secret=?", keySecret).
		And("status=?",1).Get(&s.apps)
	if err!=nil{
		return token,err
	}
	if !has{
		return token, errors.New("该应用尚未注册")
	}

	//获取token
	build:=[]byte(fmt.Sprintf("%s.%s.%s",keyId,time.Now(),keySecret))
	token=fmt.Sprintf("%x",sha256.Sum256(build))
	s.apps.Token=&token
	_, err = s.db.Cols("token").Where("id=?", s.apps.Id).Update(&s.apps)
	if err!=nil{
		return "",err
	}


	err = s.SetCache(token,s.apps.Id)
	if err!=nil{
		log.Fatal("app token缓存存储失败")
	}
	return
}

//设置app token缓存
func (s *AppService) SetCache(token string,appId uint) error{
	cacheKey:=fmt.Sprintf("%d_%s",cache.AppsTokensMap,token)
	return s.redis.Set(cacheKey, appId, time.Hour*24*7).Err()
}

//获取app token缓存
func (s *AppService) GetCache(token string) (string,error){
	cacheKey:=fmt.Sprintf("%d_%s",cache.AppsTokensMap,token)
	return s.redis.Get(cacheKey).Result()
}