package event

import (
	"fmt"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"im/model"
	"im/service"
	"im/service/orm"
	"strconv"
)

func Chat() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
		var(
			userMsg=&model.Messages{}
			err=msg.Unmarshal(&userMsg)
			result int
			ctx = websocket.GetContext(nsConn.Conn)
			user =ctx.Values().Get("user").(model.Users)
			receipt=model.Receipts{Mid:userMsg.Mid,Type:dao.SENT}
			db=orm.GetDB()
		)

		if err!=nil{
			return err
		}

		result,err=service.NewUsersUsersService().EachOtherFriends(user.AppsId,userMsg.To,userMsg.From)
		if result!=service.EACH_OTHER_FRIENDS_IS_TRUE{
			//校验失败时返回消息
			//userMsg.Err=errors.New(common.CONTACTS_BROKEN)
			//fmt.Println(userMsg)
			//str,_=userMsg.Marshal()
			msg.Err=common.NewErrorRes(common.CONTACTS_BROKEN,"你不在对方好友列表",userMsg)
			goto ERROR
		}


		if userMsg.Status==1{
			userMsg.AppsId=user.AppsId
			_,err=db.InsertOne(userMsg)
			if err!=nil{

				msg.Err=common.NewErrorRes(common.SQL_INSERT_FAILD,"消息存储失败",userMsg)
				nsConn.Emit("chat",websocket.Marshal(msg))
				return nil
			}

			//str,_=userMsg.Marshal()
			msg.Body,msg.To=websocket.Marshal(userMsg),strconv.FormatUint(userMsg.To,10)

			fmt.Println(string(msg.Serialize()))

			//TODO act报文反馈 已发送 已阅读
			dao.ChatReceiptsDao().Add(userMsg.Mid,userMsg.From,userMsg.To)

			nsConn.Emit("receipt",receipt.Marshal())

			nsConn.Conn.Server().Broadcast(nil,msg)
			return nil

		//撤回消息
		}else{
			if userMsg.From!=user.Id{
				msg.Err=common.NewErrorRes(common.ATTRIBUTION_ERROR,"无权限操作",userMsg)
				goto ERROR
			}
			_,err=dao.NewMessagesDao().UpdateByUidMid(userMsg.Mid,userMsg.From,userMsg,"status")
			if err!=nil{
				msg.Err=common.NewErrorRes(common.SQL_UPDATE_FAILD,"消息撤销失败:"+err.Error(),userMsg)
				goto ERROR
			}

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

