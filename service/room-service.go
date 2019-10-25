package service

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type RoomService interface {
	//GetAll() []model.Chatrooms
	//GetByID(id int64) (model.Chatrooms, bool)
	DeleteByID(apps_id uint,id uint64,uid uint64) (int64, error)
	Create(data *model.Chatrooms) (int64, error)
}

func NewRoomService() RoomService {
	return &roomService{db:orm.GetDB()}
}

type roomService struct {
	db *xorm.Engine
}

//func (r roomService) GetAll() []model.Chatrooms {
//
//}

//func (r roomService) GetByID(id int64) (model.Chatrooms, bool) {
//	panic("implement me")
//}

func (r roomService) DeleteByID(apps_id uint,id uint64,uid uint64) (int64,error) {
	m:=new(model.Chatrooms)
	m.Status=0
	return r.db.Cols("status").Where("id=?",id).Where("apps_id=?",apps_id).Where("uid=?",uid).Update(m)
}

func (r roomService) Create(data *model.Chatrooms) (int64, error) {
	return r.db.InsertOne(data)
}

