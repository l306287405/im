package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type appsDao struct {
	db *xorm.Engine
}

func NewAppsDao() *appsDao{
	return &appsDao{db:orm.GetDB()}
}

func (d *appsDao) GetInfoByToken(token string) *model.Apps{
	m:=new(model.Apps)
	has,err:=d.db.Where("token=?",token).Where("status=?",1).Get(m)
	if err!=nil{
		return nil
	}
	if !has{
		return nil
	}
	return m
}