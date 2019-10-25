package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type ChatroomsUsersDao struct {
	db *xorm.Engine
}

func NewChatroomsUsersDao() *ChatroomsUsersDao{
	return &ChatroomsUsersDao{db:orm.GetDB()}
}

func (d *ChatroomsUsersDao) Create(users model.ChatroomsUsers) (int64,error){
	return d.db.InsertOne(users)
}

func (d *ChatroomsUsersDao) RelationExist(appId uint,roomId uint64,userId uint64) bool{
	m:=new(model.ChatroomsUsers)
	_,err:=d.db.Cols("apps_id").Where("apps_id=?",appId).Where("room_id=?",roomId).
		Where("uid=?",userId).Get(m)
	if err!=nil{
		panic(err)
	}
	if m.AppsId>0{
		return true
	}
	return false
}

func (d *ChatroomsUsersDao) Update(appId uint,roomId uint64,userId uint64,data *model.ChatroomsUsers,cols ...string) (int64,error){
	return d.db.Cols(cols...).Where("apps_id=?",appId).Where("room_id=?",roomId).Where("uid=?",userId).Update(data)
}

func (d *ChatroomsUsersDao) IsManager(appId uint,roomId uint64,userId uint64) bool{
	chatroomsUsers:=new(model.ChatroomsUsers)
	_,err:=d.db.Cols("status").Where("apps_id=?",appId).Where("room_id=?",roomId).
		Where("uid=?",userId).Get(chatroomsUsers)
	if err!=nil{
		panic(err)
	}
	if chatroomsUsers.Status == 0 || chatroomsUsers.Status==1{
		return true
	}
	return false
}