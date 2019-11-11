package dao

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
	"log"
)

type AttachmentsDao struct {
	db *xorm.Engine
}

func NewAttachmentsDao() *AttachmentsDao{
	return &AttachmentsDao{db:orm.GetDB()}
}

func (d *AttachmentsDao) GetFileBySha1(app_id uint,sha1 string) *model.Attachments{
	m:=&model.Attachments{}
	ok,err:=d.db.Where("sha1=?",sha1).Get(m)
	if err!=nil{
		log.Println("Attachments err:"+err.Error())
		return nil
	}
	if !ok{
		return nil
	}
	return m
}

func (d *AttachmentsDao) Create(m *model.Attachments) (int64,error){
	return d.db.InsertOne(m)
}