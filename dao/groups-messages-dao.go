package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type GroupsMessagesDao struct {
	db *xorm.Engine
}

func NewGroupsMessagesDao() *GroupsMessagesDao{
	return &GroupsMessagesDao{db:orm.GetDB()}
}

func (d *GroupsMessagesDao) UpdateByUidMid(mid uint64,from uint64,data *model.GroupsMessages,cols ...string ) (int64,error){
	return d.db.Cols(cols...).Where("mid=?",mid).Where("from=?",from).Update(data)
}