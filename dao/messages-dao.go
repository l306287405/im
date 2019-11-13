package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type MessagesDao struct {
	db *xorm.Engine
}

func NewMessagesDao() *MessagesDao{
	return &MessagesDao{db:orm.GetDB()}
}

func (d *MessagesDao) UpdateByUidMid(mid uint64,from uint64,data *model.Messages,cols ...string ) (int64,error){
	return d.db.Cols(cols...).Where("mid=?",mid).Where("from=?",from).Update(data)
}