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
			num int64
			str []byte
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
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_INSERT_FAILD,err.Error()
				goto ERROR
			}

			//TODO act报文反馈 已发送 已阅读
			str,_=userMsg.Marshal()
			msg.Body=str
			fmt.Println(string(msg.Serialize()))
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil

		//撤回消息
		}else{
			if userMsg.From!=user.Id{
				userMsg.ErrCode,userMsg.ErrMsg=new(string),new(string)
				*userMsg.ErrCode,*userMsg.ErrMsg=common.SQL_INSERT_FAILD,"撤销失败,非法操作"
				goto ERROR
			}
			num,err=dao.NewGroupsMessagesDao().UpdateByUidMid(userMsg.Mid,userMsg.From,userMsg,"status")
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
		nsConn.Emit("group",str)
		return nil
	}
}
