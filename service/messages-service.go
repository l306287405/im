package service

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type MessagesService struct {
	db *xorm.Engine
}

func NewMessagesService() *MessagesService{
	return &MessagesService{db:orm.GetDB()}
}

func (s *MessagesService) GetList(appId uint,to uint64,beginTime string,endTime *string,from *uint64,limit *int,cursor *int) (*[]model.Messages,*int){
	var(
		list=new([]model.Messages)
		err error
	)

	session:=s.db.Where("apps_id=? and `to`=? and create_at>=?",appId,to,beginTime)
	if from!=nil{
		session=session.Where("from=?",*from)
	}

	if endTime!=nil{
		session=session.Where("create_at<=?",*endTime)
	}

	if limit!=nil{
		*limit+=1
		if cursor==nil{
			cursor=new(int)
			*cursor=0
		}
		session=session.Limit(*limit,*cursor)
	}
	err=session.Desc("create_at").Find(list)
	if err!=nil{
		panic(err)
	}
	if limit!=nil && len(*list)>*limit{
		*cursor+=*limit
	}
	return list,cursor
}