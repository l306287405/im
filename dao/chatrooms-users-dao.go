package dao

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

type ChatroomsUsersDao struct {
	db *xorm.Engine
	cache *redis.Client
}

func NewChatroomsUsersDao() *ChatroomsUsersDao{
	return &ChatroomsUsersDao{db:orm.GetDB(),cache:cache.Init()}
}

func (d *ChatroomsUsersDao) Create(users model.ChatroomsUsers) (int64,error){
	return d.db.InsertOne(users)
}

func (d *ChatroomsUsersDao) RelationExist(appId uint,roomId uint64,userId uint64) *int8{
	m:=new(model.ChatroomsUsers)
	status,err:=d.cache.Get(d.getCacheKey(appId,roomId,userId)).Int()
	if err==nil{
		tempStatus:=int8(status)
		return &tempStatus
	}

	ok,err:=d.db.Cols("status").Where("apps_id=?",appId).Where("room_id=?",roomId).
		Where("uid=?",userId).Get(m)
	if err!=nil{
		panic(err)
	}
	if !ok{
		err=d.setCache(d.getCacheKey(appId,roomId,userId),-1)
		if err!=nil{
			log.Println("群老与用户关系缓存失败 原因:"+err.Error())
		}
		return nil
	}
	err=d.setCache(d.getCacheKey(appId,roomId,userId),m.Status)
	if err!=nil{
		log.Println("群老与用户关系缓存失败 原因:"+err.Error())
	}
	return &m.Status
}

func (d *ChatroomsUsersDao) getCacheKey(appId uint,roomId uint64,userId uint64) string{
	return fmt.Sprintf("%d_%d_%d_%d",cache.ROOMS_USERS_MAP,appId,roomId,userId)
}

func (d *ChatroomsUsersDao) setCache(cacheKey string,status int8) error{
	return d.cache.Set(cacheKey,status,time.Hour*24).Err()
}

func (d *ChatroomsUsersDao) Update(appId uint,roomId uint64,userId uint64,data *model.ChatroomsUsers,cols ...string) (int64,error){
	return d.db.Cols(cols...).Where("apps_id=?",appId).Where("room_id=?",roomId).Where("uid=?",userId).Update(data)
}

func (d *ChatroomsUsersDao) IsManager(appId uint,roomId uint64,userId uint64) bool{
	chatroomsUsers:=new(model.ChatroomsUsers)
	find,err:=d.db.Cols("status").Where("apps_id=?",appId).Where("room_id=?",roomId).
		Where("uid=?",userId).Get(chatroomsUsers)
	if err!=nil{
		panic(err)
	}
	if !find{
		return false
	}
	if chatroomsUsers.Status == 0 || chatroomsUsers.Status==1{
		return true
	}
	return false
}