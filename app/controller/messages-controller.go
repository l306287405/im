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

type messagesGetParams struct {
	BeginTime string  `json:"begin_time"`
	EndTime   *string `json:"end_time,omitempty"`
	From      *uint64 `json:"from"`
	Limit     *int    `json:"limit,omitempty"`
	Cursor    *int    `json:"cursor,omitempty"`
}

func (c *MessagesController) Get() {
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		response=new(messagesGetResponse)

		//params
		params=&messagesGetParams{}
		beginTime=c.Ctx.URLParamTrim("begin_time")
		endTime=c.Ctx.URLParamTrim("end_time")
		from=uint64(c.Ctx.URLParamInt64Default("from",0))
		limit=c.Ctx.URLParamIntDefault("limit",0)
		cursor=c.Ctx.URLParamIntDefault("cursor",0)

		err error
	)

	if beginTime == "" {
		err=errors.New("参数缺失或错误")
		goto PARAMS_ERR
	}
	params.BeginTime=beginTime
	if endTime!=""{
		params.EndTime=&endTime
	}
	if from!=0{
		params.From=&from
	}
	if limit!=0{
		params.Limit=&limit
	}
	if cursor!=0{
		params.Cursor=&cursor
	}

	response.List,response.Cursor=service.NewMessagesService().GetList(user.AppsId,user.Id,params.BeginTime,params.EndTime,params.From,params.Limit,params.Cursor)
	c.Ctx.JSON(common.SendSmile(response))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return


}