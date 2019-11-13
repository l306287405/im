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
			num int64
			str []byte
			receipt=model.Receipts{Mid:userMsg.Mid,Type:model.SENT}
			db=orm.GetDB()
		)

		if err!=nil{
			return err
		}

		result,err=service.NewUsersUsersService().EachOtherFriends(user.AppsId,userMsg.To,userMsg.From)
		if result!=service.EACH_OTHER_FRIENDS_IS_TRUE{
			//校验失败时返回消息
			userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
			*userMsg.ErrCode,*userMsg.ErrMsg=common.CONTACTS_BROKEN,"消息发送失败,您不在对方好友列表."

			str,_=userMsg.Marshal()
			nsConn.Emit("chat",str)
			return nil
		}


		if userMsg.Status==1{
			userMsg.AppsId=user.AppsId
			_,err=db.InsertOne(userMsg)
			if err!=nil{
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_INSERT_FAILD,"消息存储失败"
				str,_=userMsg.Marshal()
				nsConn.Emit("chat",str)
				return nil
			}

			str,_=userMsg.Marshal()
			msg.Body,msg.To=str,strconv.FormatUint(userMsg.To,10)
			fmt.Println(string(msg.Serialize()))

			//TODO act报文反馈 已发送 已阅读
			model.ChatReceipts().Add(userMsg.Mid,userMsg.From,userMsg.To)

			nsConn.Emit("receipt",receipt.Marshal())

			nsConn.Conn.Server().Broadcast(nil,msg)
			return nil

		}else{
			if userMsg.From!=user.Id{
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_INSERT_FAILD,"撤销失败,非法操作"
				goto ERROR
			}
			num,err=dao.NewMessagesDao().UpdateByUidMid(userMsg.Mid,userMsg.From,userMsg,"status")
			if err!=nil{
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_ERROR,err.Error()
				goto ERROR
			}

			//无状态改变
			if num==0{
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_UPDATE_FAILD,"重复提交请求,或者状态已被变更"
				goto ERROR
			}

			str,_=userMsg.Marshal()
			msg.Body=str
			fmt.Println(string(msg.Serialize()))
			nsConn.Conn.Server().Broadcast(nsConn, msg)


			return nil
		}

	ERROR:
		str,_=userMsg.Marshal()
		nsConn.Emit("chat",str)
		return nil
	}
}

