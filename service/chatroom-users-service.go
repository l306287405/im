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
	GetListByUser(appsId uint,uid uint64) (*[]uint64,error)
}

func NewChatroomUsersService() ChatroomUsersService {
	return &chatroomUsersService{db:orm.GetDB()}
}

type chatroomUsersService struct {
	db *xorm.Engine
}

func (c *chatroomUsersService) GetListByUser(appsId uint, uid uint64) (r *[]uint64, e error) {
	list:=&[]model.ChatroomsUsers{}
	err:=c.db.Cols("room_id").Where("apps_id=?",appsId).Where("uid=?",uid).Where("status=?",1).Find(list)
	if err!=nil{
		return nil,err
	}
	r=new([]uint64)
	for _,v := range *list{
		*r = append(*r,v.RoomId)
	}
	return r,nil
}

func (c *chatroomUsersService) DeleteByID(appsId uint, roomId uint64, uid uint64) (int64, error) {
	data:=&model.ChatroomsUsers{Status:-1}
	return dao.NewChatroomsUsersDao().Update(appsId,roomId,uid,data,"status")
}

func (c *chatroomUsersService) Create(data *model.ChatroomsUsers) (int64, error) {
	d:=dao.NewChatroomsUsersDao()
	exist:=d.RelationExist(data.AppsId,data.RoomId,data.Uid)
	if exist!=nil{
		d.DeleteCache(d.GetCacheKey(data.AppsId,data.RoomId,data.Uid))
		return c.db.Cols("role","status","Joined_at").Where("apps_id=?",data.AppsId).Where("room_id=?",data.RoomId).Where("uid=?",data.Uid).
			Where("role!=?",0).Update(data)
	}
	return c.db.InsertOne(data)
}


