package dao

import (
	"github.com/go-xorm/xorm"
	"im/service/orm"
)

type ChatroomsDao struct {
	db *xorm.Engine
}

func NewChatroomsDao() *ChatroomsDao{
	return &ChatroomsDao{db:orm.GetDB()}
}