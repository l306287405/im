package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type CollectionsDao struct {
	db *xorm.Engine
}

func NewCollectionsDao() *CollectionsDao{
	return &CollectionsDao{db:orm.GetDB()}
}

func (d *CollectionsDao) Create(m *model.Collections) (int64,error){
	return d.db.InsertOne(m)
}

func (d *CollectionsDao) GetList(appId uint,userId uint64,textType *byte) (*[]model.Collections,error){
	session:=d.db.Where("apps_id=?",appId).Where("uid=?",userId)
	if textType!=nil{
		session=session.Where("text_type=?",*textType)
	}
	list:=&[]model.Collections{}
	err:=session.Find(list)
	return list,err
}

func (d *CollectionsDao) Delete(id uint64,appId uint,userId uint64) (int64,error){
	m:=&model.Collections{}
	return d.db.Where("id=?",id).Where("apps_id=?",appId).Where("uid=?",userId).Delete(m)
}