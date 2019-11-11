package controller

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/dao"
	"im/model"
	"net/http"
)

type CollectionsController struct {
	Session *sessions.Session
	Ctx iris.Context
}

type collectionsGetParams struct {
	Limit uint		`json:"limit"`
	Offset uint64	`json:"offset"`
	TextType *byte	`json:"text_type"`
}

func (c *CollectionsController) Get(){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		params=new(collectionsGetParams)
		list=&[]model.Collections{}
		err=c.Ctx.ReadJSON(params)
	)

	if err!=nil{
		goto PARAMS_ERR
	}

	if params.Limit<1 || params.Limit>20 {
		err=errors.New("limit must between 1 and 20")
		goto PARAMS_ERR
	}

	list,err=dao.NewCollectionsDao().GetList(user.AppsId,user.Id,params.TextType)
	if err!=nil{
		goto SQL_ERR
	}

	c.Ctx.JSON(common.SendSmile(list))

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

}

func (c *CollectionsController) Post(){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		params=&model.Collections{}
		err=c.Ctx.ReadJSON(params)
	)

	if err!=nil{
		goto PARAMS_ERR
	}

	params.AppsId,params.Uid=user.AppsId,user.Id

	_,err=dao.NewCollectionsDao().Create(params)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(params.Id,"收藏成功"))
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

func (c *CollectionsController) DeleteBy(id uint64){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		num int64
		err error
	)

	num,err=dao.NewCollectionsDao().Delete(id,user.AppsId,user.Id)
	if err!=nil{
		goto SQL_ERR
	}

	c.Ctx.JSON(common.SendSmile(num,"删除收藏成功"))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

}