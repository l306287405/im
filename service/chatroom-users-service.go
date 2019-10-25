package service

import (
	"github.com/go-xorm/xorm"
	"im/dao"
	"im/model"
	"im/service/orm"
)

type ChatroomUsersService interface {
	//GetAll() []model.Chatrooms
	//GetByID(id int64) (model.Chatrooms, bool)
	DeleteByID(appsId uint,roomId uint64,uid uint64) (int64, error)
	Create(data *model.ChatroomsUsers) (int64, error)
}

func NewChatroomUsersService() ChatroomUsersService {
	return &chatroomUsersService{db:orm.GetDB()}
}

type chatroomUsersService struct {
	db *xorm.Engine
}

func (c *chatroomUsersService) DeleteByID(appsId uint, roomId uint64, uid uint64) (int64, error) {
	data:=&model.ChatroomsUsers{Status:-1}
	return dao.NewChatroomsUsersDao().Update(appsId,roomId,uid,data,"status")
}

func (c *chatroomUsersService) Create(data *model.ChatroomsUsers) (int64, error) {
	d:=dao.NewChatroomsUsersDao()
	exist:=d.RelationExist(data.AppsId,data.RoomId,data.Uid)
	if exist{
		return d.Update(data.AppsId,data.RoomId,data.Uid,data,"role","status","Joined_at")
	}
	return c.db.InsertOne(data)
}


