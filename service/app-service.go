package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"im/dao"
	"im/model"
	"im/service/cache"
	"im/service/orm"
	"strconv"
	"time"
)

type AppService struct {
	apps model.Apps
	db *xorm.Engine
	redis *redis.Client
}

func NewAppService() *AppService{
	return &AppService{db:orm.GetDB(),redis:cache.Init()}
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

	return
}

//设置app token缓存
func (s *AppService) SetCache(token string,appId uint) error{

	tokenCacheKey:=fmt.Sprintf("%d_%s",cache.APPS_TOKEN_MAP,token)
	IdCacheKey:=fmt.Sprintf("%d_%d",cache.APPS_TOKEN_MAP,appId)

	//清理旧缓存
	lastToken,err := s.redis.Get(IdCacheKey).Result()
	if err==nil{
		lastTokenCacheKey:=fmt.Sprintf("%d_%s",cache.APPS_TOKEN_MAP,lastToken)
		s.redis.Del(lastTokenCacheKey)
	}

	//用于判断token是否有效的缓存
	err =s.redis.Set(tokenCacheKey, appId, 0).Err()
	if err!=nil{
		return err
	}

	//获取id当前关联token的缓存,多用于
	return s.redis.Set(IdCacheKey, token, 0).Err()
}

//获取app token缓存
func (s *AppService) GetToken(token string) (uint,error){
	cacheKey:=fmt.Sprintf("%d_%s",cache.APPS_TOKEN_MAP,token)
	str,err:=s.redis.Get(cacheKey).Result()
	if err!=nil || str==""{
		m:=dao.NewAppsDao().GetInfoByToken(token)
		if m==nil{
			return 0,err
		}

		_=s.SetCache(*m.Token,m.Id)
		return m.Id,nil
	}
	appId,err:=strconv.ParseUint(str,10,32)
	return uint(appId),err
}