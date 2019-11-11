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

type MessagesController struct {
	Session *sessions.Session
	Ctx iris.Context
}

type messagesGetResponse struct {
	List *[]model.Messages
	Cursor *int
}

func (c *MessagesController) Get() {
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		response=new(messagesGetResponse)
		err error

		//params
		beginTime=c.Ctx.URLParamTrim("begin_time")
		endTime *string
		from *uint64
		limit *int
		cursor *int
	)

	if beginTime == "" {
		err=errors.New("参数缺失或错误")
		goto PARAMS_ERR
	}

	if c.Ctx.URLParamTrim("end_time")!=""{
		endTime=new(string)
		*endTime=c.Ctx.URLParamTrim("end_time")
	}

	if f,err:=c.Ctx.URLParamInt64("from");err==nil{
		from=new(uint64)
		*from=uint64(f)
	}

	if l,err:=c.Ctx.URLParamInt("limit");err==nil{
		limit=new(int)
		*limit=l
	}

	if c,err:=c.Ctx.URLParamInt("cursor");err==nil{
		cursor=new(int)
		*cursor=c
	}

	response.List,response.Cursor=service.NewMessagesService().GetList(user.AppsId,user.Id,beginTime,endTime,from,limit,cursor)
	c.Ctx.JSON(common.SendSmile(response))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return


}