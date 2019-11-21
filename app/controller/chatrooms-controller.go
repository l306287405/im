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

type ChatroomsController struct {
	Session *sessions.Session
	Ctx iris.Context
}

func (c *ChatroomsController) Get(){
	var(
		user=c.Ctx.Values().Get("user").(model.Users)
		list=&[]model.RelationGroupsUsers{}

		err error
	)

	list,err=dao.NewChatroomsUsersDao().GetListByUid(user.AppsId,user.Id)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(list))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return
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

	ruModel=&model.ChatroomsUsers{AppsId:user.AppsId,RoomId:room.Id,Uid:user.Id,Role:model.ROOM_USER_ROLE_IS_OWNER,Status:1,JoinedAt:new(string)}
	*ruModel.JoinedAt=common.DateTime(time.Now())
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
		num int64
		err error
	)

	user=c.Ctx.Values().Get("user").(model.Users)

	num,err=roomService.DeleteByID(user.AppsId,id,user.Id)
	if err!=nil{
		goto SQL_ERR
	}

	if num!=1{
		num, err = service.NewChatroomUsersService().DeleteByID(user.AppsId, id, user.Id)
		if num==1{
			c.Ctx.JSON(common.SendSmile(1,"群聊退出成功"))
		}else{
			c.Ctx.JSON(common.SendSad("群聊退出失败,非群内成员或重复操作"))
		}
	}else{
		_,err = service.NewChatroomUsersService().DeleteByRoomId(user.AppsId,id)
		c.Ctx.JSON(common.SendSmile(1,"群聊关闭成功"))
	}
	return


SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

}

type roomsPut struct {
	Name *string `json:"name"`
	Desc *string `json:"desc"`
	MaxUsers *uint16 `json:"max_users"`
	Approval *byte `json:"approval"`
	Status *byte `json:"status"`
}
func (c *ChatroomsController) PutBy(roomId uint64){
	var(
		user model.Users
		chatrooms =new(model.Chatrooms)
		roomService=service.NewRoomService()
		putParams=&roomsPut{}
		changeStr []string
		err error
	)

	user=c.Ctx.Values().Get("user").(model.Users)
	if !dao.NewChatroomsUsersDao().IsManager(user.AppsId,roomId,user.Id){
		err=errors.New("非法修改群聊房间")
		goto PARAMS_ERR
	}

	err=c.Ctx.ReadJSON(putParams)
	if err!=nil{
		goto PARAMS_ERR
	}
	if status:=putParams.Status;status!=nil{
		if *status==0 || *status==1{
			changeStr=append(changeStr,"status")
			chatrooms.Status=*status
		}
	}

	if name:=putParams.Name;name!=nil{
		changeStr=append(changeStr,"name")
		chatrooms.Name=*name
	}

	if desc:=putParams.Desc;desc!=nil{
		changeStr=append(changeStr,"desc")
		chatrooms.Desc=*desc
	}

	if maxUsers:=putParams.MaxUsers;maxUsers!=nil{
		if *maxUsers >=200 && *maxUsers <=2000{
			changeStr=append(changeStr,"max_users")
			chatrooms.MaxUsers=maxUsers
		}
	}

	if approval:=putParams.Approval;approval!=nil{
		if *approval==0 || *approval==1{
			changeStr=append(changeStr,"approval")
			chatrooms.Approval=approval
		}
	}

	err=roomService.UpdateById(user.AppsId,roomId,chatrooms,changeStr...)
	if err!=nil{
		goto SQL_ERR
	}
	c.Ctx.JSON(common.SendSmile(1,"修改成功"))
	return

PARAMS_ERR:
	c.Ctx.StatusCode(http.StatusBadRequest)
	c.Ctx.JSON(common.SendCry("参数错误:"+err.Error()))
	return

SQL_ERR:
	c.Ctx.StatusCode(http.StatusInternalServerError)
	c.Ctx.JSON(common.SendSad("服务器发生错误 "+err.Error()))
	return

}