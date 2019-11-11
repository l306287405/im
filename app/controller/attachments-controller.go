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

type AttachmentsController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c *AttachmentsController) Post(){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		m=&model.Attachments{}
		err=c.Ctx.ReadJSON(m)
	)

	if err!=nil{
		goto PARAMS_ERR
	}

	if m.Url==""{
		err=errors.New("资源路径缺失")
		goto PARAMS_ERR
	}
	if m.Mime==""{
		err=errors.New("mime类型缺失")
		goto PARAMS_ERR
	}
	if m.Size==0{
		err=errors.New("文件尺寸缺失")
		goto PARAMS_ERR
	}
	if m.Sha1==""{
		err=errors.New("sha1值缺失")
		goto PARAMS_ERR
	}

	m.AppsId=user.AppsId
	m.Status=1

	_,err=dao.NewAttachmentsDao().Create(m)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(m.Id))
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

func (c *AttachmentsController) GetBy(sha1 string){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		m=&model.Attachments{}
	)

	m=dao.NewAttachmentsDao().GetFileBySha1(user.AppsId,sha1)

	if m==nil{
		c.Ctx.StatusCode(http.StatusNotFound)
		c.Ctx.JSON(common.SendSad("指定附件资源不存在你"))
		return
	}
	c.Ctx.JSON(common.SendSmile(m))
}