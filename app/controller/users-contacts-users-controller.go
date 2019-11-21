package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/model"
	"im/service"
	"im/service/orm"
	"net/http"
	"strconv"
)

type UsersContactsUsersController struct {
	Session *sessions.Session
	Ctx iris.Context
}


//创建新联系人
func (c *UsersContactsUsersController) Post(){
	var(
		target model.Users
		uu model.UsersUsers
		db=orm.GetDB()
		fid=c.Ctx.Params().Get("fid")
		user=c.Ctx.Values().Get("user").(model.Users)
		found bool
		cid uint64
		id int64

		err error
	)

	target.Id,err=strconv.ParseUint(fid,10,64)

	if err!=nil{
		goto PARAMS_ERR
	}

	found,err=db.Where("id=?",target.Id).Where("apps_id=?",user.AppsId).Get(&target)
	if err!=nil{
		goto SQL_ERR
	}
	if !found{
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.JSON(common.SendSad("指定用户不存在"))
		return
	}

	uu=model.UsersUsers{AppsId:user.AppsId,Uid:user.Id,Fid:target.Id,Cid:cid}
	id,err=db.InsertOne(uu)

	if err!=nil{
		goto SQL_ERR
	}

	err=service.NewUsersUsersService().DelCacheOfEOF(user.AppsId,user.Id,target.Id)
	c.Ctx.JSON(common.SendSmile(id,"添加好友成功"))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("参数错误"))
	return
}

//删除好友关系
func (c *UsersContactsUsersController) Delete(){
	var(
		target model.Users
		uu model.UsersUsers
		db=orm.GetDB()
		err error
		fid=c.Ctx.Params().Get("fid")
		user model.Users
	)

	target.Id,err=strconv.ParseUint(fid,10,64)

	if err!=nil{
		goto PARAMS_ERR
	}

	user=c.Ctx.Values().Get("user").(model.Users)
	_,err=db.Where("uid=?",user.Id).Where("fid=?",target.Id).Delete(&uu)
	if err!=nil{
		goto SQL_ERR
	}

	err=service.NewUsersUsersService().DelCacheOfEOF(user.AppsId,user.Id,target.Id)
	if err!=nil{
		goto SQL_ERR
	}

	c.Ctx.JSON(common.SendSmile(1,"删除好友成功"))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("参数错误"))
	return

}