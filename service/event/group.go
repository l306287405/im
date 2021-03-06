package event

import (
	"fmt"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"im/model"
	"im/service/orm"
)

func Group() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
		var(
			userMsg=&model.GroupsMessages{}
			ctx = websocket.GetContext(nsConn.Conn)
			userInter=ctx.Values().Get("user")
			user=model.Users{}
			err error
			db=orm.GetDB()
		)
		if userInter==nil{
			fmt.Println("用户信息有误,请检查")
			return nil
		}
		user=userInter.(model.Users)

		err=msg.Unmarshal(&userMsg)
		if err!=nil{
			return err
		}

		//发送消息
		if userMsg.Status==1{
			//插入消息记录
			userMsg.AppsId=user.AppsId
			_,err=db.InsertOne(userMsg)
			if err!=nil{
				msg.Err=common.NewErrorRes(common.SQL_INSERT_FAILD,err.Error(),userMsg)
				goto ERROR
			}

			//TODO act报文反馈 已发送 已阅读
			//str,_=userMsg.Marshal()
			msg.Body=websocket.Marshal(userMsg)
			fmt.Println(string(msg.Serialize()))
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil

		//撤回消息
		}else{
			if userMsg.From!=user.Id{
				msg.Err=common.NewErrorRes(common.ATTRIBUTION_ERROR,"无权限操作",userMsg)
				goto ERROR
			}
			_,err=dao.NewGroupsMessagesDao().UpdateByUidMid(userMsg.Mid,userMsg.From,userMsg,"status")
			if err!=nil{
				msg.Err=common.NewErrorRes(common.SQL_UPDATE_FAILD,"消息撤销失败:"+err.Error(),userMsg)
				goto ERROR
			}

			//str,_=userMsg.Marshal()
			msg.Body=websocket.Marshal(userMsg)
			fmt.Println(string(msg.Serialize()))
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		}

ERROR:
		nsConn.Conn.Write(msg)
		return nil
	}
}
