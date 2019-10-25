package controller

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"im/common"
	"im/dao"
	"im/model"
	"im/service"
	"net/http"
	"time"
)

type ChatroomsUsersController struct {
	Session *sessions.Session
	Ctx iris.Context
}

type postParams struct {
	UserId uint64 `json:"user_id"`
	Status int8	  `json:"status"`
}

func (c *ChatroomsUsersController) Post(){
	var(
		user         =c.Ctx.Values().Get("user").(model.Users)
		roomId,err   =c.Ctx.Params().GetUint64("room_id")
		params		 =new(postParams)
		usersService =service.NewChatroomUsersService()
		data         =new(model.ChatroomsUsers)
		changed		 int64
		joined_at 	 string
	)
	if err!=nil{
		goto PARAMS_ERR
	}
	err=c.Ctx.ReadJSON(params)
	if err!=nil{
		goto PARAMS_ERR
	}

	if params.UserId==0{
		err=errors.New("用户id缺失或无效")
		goto PARAMS_ERR
	}

	//确认房间以及归属权
	if !dao.NewChatroomsUsersDao().IsManager(user.AppsId,roomId,user.Id) {
		err=errors.New("非法操作")
	}

	if params.Status==1{
		joined_at=time.Now().Format("2006-01-02 15:04:05")
	}

	data=&model.ChatroomsUsers{Uid:params.UserId,AppsId:user.AppsId,RoomId:roomId,Role:model.ROOM_USER_ROLE_IS_MEMBER,
		Status:params.Status,JoinedAt:&joined_at}
	changed,err= usersService.Create(data)

	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(changed,"请求成功"))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return

}

type deleteParams struct {
	UserId uint64 `json:"user_id"`
}

func (c *ChatroomsUsersController) Delete(){
	var(
		user         =c.Ctx.Values().Get("user").(model.Users)
		roomId,err   =c.Ctx.Params().GetUint64("room_id")
		params    	 =new(deleteParams)
		usersService =service.NewChatroomUsersService()
		changed		 int64
	)

	if err!=nil{
		goto PARAMS_ERR
	}

	err=c.Ctx.ReadJSON(params)
	if err!=nil{
		goto PARAMS_ERR
	}

	//确认房间以及归属权
	if !dao.NewChatroomsUsersDao().IsManager(user.AppsId,roomId,user.Id) || params.UserId!=user.Id {
		err=errors.New("非法操作")
	}

	changed,err=usersService.DeleteByID(user.AppsId,roomId,params.UserId)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(changed,"退出群聊成功"))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return


PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("错误 "+err.Error()))
	return

}