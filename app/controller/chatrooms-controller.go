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

type ChatroomsController struct {
	Session *sessions.Session
	Ctx iris.Context
}

//群聊创建
func (c *ChatroomsController) Post(){
	var(
		room =new(model.Chatrooms)
		user model.Users
		ruModel =new(model.ChatroomsUsers)
		roomService=service.NewRoomService()
		ruService=service.NewChatroomUsersService()

		err error
	)

	if err=c.Ctx.ReadJSON(&room);err!=nil{
		goto PARAMS_ERR
	}

	if room.Name==""{
		err=errors.New("群聊名称参数缺失")
		goto PARAMS_ERR
	}
	if room.Desc==""{
		err=errors.New("群聊简介参数缺失")
		goto PARAMS_ERR
	}

	if room.MaxUsers==nil{
		*room.MaxUsers=200
	}
	if *room.MaxUsers<200 || *room.MaxUsers>2000{
		err=errors.New("最大人数不得小于200或大于2000")
		goto PARAMS_ERR
	}
	if room.Approval==nil{
		*room.Approval=0
	}
	if *room.Approval!=0 && *room.Approval!=1{
		err=errors.New("入群批准入参错误")
		goto PARAMS_ERR
	}

	user=c.Ctx.Values().Get("user").(model.Users)
	room.Uid,room.AppsId,room.Status=&user.Id,user.AppsId,1
	_,err=roomService.Create(room)
	if err!=nil{
		goto SQL_ERR
	}

	ruModel=&model.ChatroomsUsers{AppsId:user.AppsId,RoomId:room.Id,Uid:user.Id,Role:model.ROOM_USER_ROLE_IS_OWNER,Status:1}
	_,err=ruService.Create(ruModel)
	if err!=nil{
		goto SQL_ERR
	}

	c.Ctx.JSON(common.SendSmile(room.Id,"创建群聊成功"))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("参数错误 "+err.Error()))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

}

//群聊关闭
func (c *ChatroomsController) DeleteBy(id uint64){
	var(
		user model.Users
		roomService=service.NewRoomService()

		err error
	)

	user=c.Ctx.Values().Get("user").(model.Users)

	//TODO 应当要移除房间和人员的关系

	_,err=roomService.DeleteByID(user.AppsId,id,user.Id)
	if err!=nil{
		goto SQL_ERR
	}

	c.Ctx.JSON(common.SendSmile(1,"群聊关闭成功"))
	return


SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

//PARAMS_ERR:
//	c.Ctx.StatusCode(http.StatusBadRequest)
//	c.Ctx.JSON(common.SendCry("参数错误"))
//	return
}