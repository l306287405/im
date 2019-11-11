package service

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/orm"
)

type GroupsMessagesService struct {
	db *xorm.Engine
}

func NewGroupsMessagesService() *GroupsMessagesService{
	return &GroupsMessagesService{db:orm.GetDB()}
}

func (s *GroupsMessagesService) GetList(appId uint,to []uint64,beginTime string,endTime *string,limit *int,cursor *int) (*[]model.GroupsMessages,*int){
	var(
		list=new([]model.GroupsMessages)
		session=&xorm.Session{}
		err error
	)

	session=s.db.Where("apps_id=?",appId)

	if len(to)>1{
		session=session.In("`to`",to)
	}else{
		session=session.Where("`to`=?",to[0])
	}
	session=session.Where("create_at>=?",beginTime)


	if endTime!=nil{
		session=session.Where("create_at<=?",*endTime)
	}

	session=session.Desc("create_at")
	if limit!=nil{
		if cursor==nil{
			cursor=new(int)
			*cursor=0
		}
		session=session.Limit(*limit+1,*cursor)
	}
	err=session.Find(list)
	if err!=nil{
		panic(err)
	}
	if limit!=nil{
		if len(*list)>*limit{
			*cursor+=*limit
			*list=(*list)[:len(*list)-1]
		}else{
			*cursor=0
		}
	}
	return list,cursor
}