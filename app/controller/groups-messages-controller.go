package controller

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/model"
	"im/service"
	"net/http"
)

type GroupsMessagesController struct {
	Session *sessions.Session
	Ctx iris.Context
}

type groupsMessagesGetResponse struct {
	List *[]model.GroupsMessages	`json:"list"`
	Cursor *int	`json:"cursor"`
}

func (c *GroupsMessagesController) Get() {
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		response=new(groupsMessagesGetResponse)
		err error

		//params
		beginTime=c.Ctx.URLParamTrim("begin_time")
		endTime *string
		to []uint64
		limit *int
		cursor *int
	)

	if c.Ctx.URLParamExists("to"){
		to=append(to,uint64(c.Ctx.URLParamInt64Default("to",0)))
	}else{
		tos,err:=service.NewChatroomUsersService().GetListByUser(user.AppsId,user.Id)
		if err!=nil{
			goto SQL_ERR
		}
		to=*tos
	}

	if beginTime == "" {
		err=errors.New("参数缺失或错误")
		goto PARAMS_ERR
	}

	if c.Ctx.URLParamExists("end_time"){
		endTime=new(string)
		*endTime=c.Ctx.URLParamTrim("end_time")
	}

	if l,err:=c.Ctx.URLParamInt("limit");err==nil{
		limit=new(int)
		*limit=l
	}

	if c,err:=c.Ctx.URLParamInt("cursor");err==nil{
		cursor=new(int)
		*cursor=c
	}

	response.List,response.Cursor=service.NewGroupsMessagesService().GetList(user.AppsId,to,beginTime,endTime,limit,cursor)
	c.Ctx.JSON(common.SendSmile(response))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return
}
