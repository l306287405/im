package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/cache"
	"im/service/orm"
	"log"
	"time"
)

type UsersUsersService struct {
	db *xorm.Engine
	cache *redis.Client
}

const EACH_OTHER_FRIENDS_IS_TRUE = 1
const EACH_OTHER_FRIENDS_IS_FALSE = 0
const EACH_OTHER_FRIENDS_IS_NIL = -1

func NewUsersUsersService() *UsersUsersService{
	return &UsersUsersService{db:orm.GetDB(),cache:cache.Init()}
}

//判断是不是互为好友
func (s *UsersUsersService) EachOtherFriends(appId uint,userA uint64,userB uint64) (result int,err error){

	defer func() {
		if err==nil{
			err=s.SetCacheOfEOF(appId,userA,userB,result)
		}
	}()

	//获取缓存中好友状态
	result,err=s.cache.Get(s.getCacheKeyOfEOF(appId,userA,userB)).Int()
	if err==nil{
		return result,nil
	}

	uu1:=new(model.UsersUsers)
	ok,err:=s.db.Where("apps_id=?",appId).Where("uid=?",userA).Where("fid=?",userB).Get(uu1)
	if err!=nil{
		return EACH_OTHER_FRIENDS_IS_NIL,err
	}
	if !ok{
		return EACH_OTHER_FRIENDS_IS_FALSE,nil
	}
	uu2:=new(model.UsersUsers)
	ok,err = s.db.Where("apps_id=?",appId).Where("uid=?",userB).Where("fid=?",userA).Get(uu2)
	if err!=nil{
		return EACH_OTHER_FRIENDS_IS_NIL,err
	}
	if !ok{
		log.Println(appId,userB,userA,uu2)
		return EACH_OTHER_FRIENDS_IS_FALSE,nil
	}
	return EACH_OTHER_FRIENDS_IS_TRUE,nil

}

func(s *UsersUsersService) getCacheKeyOfEOF(appId uint,userA uint64,userB uint64) string{
	if userA>userB{
		return fmt.Sprintf("%d_%d_%d_%d",cache.EACH_OTHER_FRIENDS,appId,userB,userA)
	}
	return  fmt.Sprintf("%d_%d_%d_%d",cache.EACH_OTHER_FRIENDS,appId,userA,userB)
}

func(s *UsersUsersService) SetCacheOfEOF(appId uint,userA uint64,userB uint64,handshake int) error{
	return s.cache.Set(s.getCacheKeyOfEOF(appId,userA,userB),handshake,time.Hour*24).Err()
}

func(s *UsersUsersService) DelCacheOfEOF(appId uint,userA uint64,userB uint64) error{
	return s.cache.Del(s.getCacheKeyOfEOF(appId,userA,userB)).Err()
}