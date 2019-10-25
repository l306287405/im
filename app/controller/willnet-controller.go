package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/neffos"
	"im/common"
	"im/dao"
	"im/model"
	"im/service"
	"log"
	"strconv"
)

type userMessage struct {
	From    uint64  `json:"from"`
	To      uint64  `json:"to"`
	Type    string  `json:"type"`
	Text    string  `json:"text"`
	ErrCode *string `json:"err_code,omitempty"`
	ErrMsg  *string	`json:"err_msg,omitempty"`
}

//在线人数,用空间换时间
var OnlineMans map[uint64]bool

func (u *userMessage) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *userMessage) Unmarshal(b []byte) error {
	return json.Unmarshal(b, u)
}

func Willnet() websocket.Events{
	return websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			user:=ctx.Values().Get("user").(model.Users)

			//加入在线人数
			err := dao.NewUsersDao().Online(user.AppsId,user.Id,nsConn.Conn.ID())
			if err!=nil{
				r,_:=json.Marshal(common.SendCry("用户上线失败,关闭链接"))
				nsConn.Emit("chatTo",r)
				nsConn.Conn.Close()
				return errors.New("用户上线失败,关闭链接")
			}

			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			log.Printf("[%s] connected to namespace [%s].",nsConn, msg.Namespace)
			return nil
		},

		//断开连接
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			user:=ctx.Values().Get("user").(model.Users)

			//从在线人数中移除
			err:= dao.NewUsersDao().OffLine(user.AppsId,user.Id)
			if err!=nil{
				r,_:=json.Marshal(common.SendCry("用户用户下线失败,关闭链接"))
				nsConn.Emit("chatTo",r)
				nsConn.Conn.Close()
				return errors.New("用户下线失败,关闭链接")
			}

			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},

		websocket.OnRoomJoin: func(nsConn *neffos.NSConn, msg neffos.Message) error {
			//TODO 校验角色跟房间的关系

			return nil
		},

		websocket.OnRoomJoined: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			user:=ctx.Values().Get("user").(model.Users)

			text := fmt.Sprintf("欢迎ID [%d] 的用户 [%s] 加入 [%s] 号房间",user.Id,user.Nickname, msg.Room)
			log.Printf("%s", text)

			// notify others.
			nsConn.Conn.Server().Broadcast(nsConn, websocket.Message{
				Namespace: msg.Namespace,
				Room:      msg.Room,
				Event:     "notify",
				Body:      []byte(text),
			})


			return nil
		},
		websocket.OnRoomLeft: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			text := fmt.Sprintf("[%s] left from room [%s].", nsConn, msg.Room)
			log.Printf("%s", text)

			// notify others.
			nsConn.Conn.Server().Broadcast(nsConn, websocket.Message{
				Namespace: msg.Namespace,
				Room:      msg.Room,
				Event:     "notify",
				Body:      []byte(text),
			})

			return nil
		},

		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			//user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)

			nsConn.Conn.Server().Broadcast(nsConn, msg)

			return nil
		},
		"chatTo": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			var(
				userMsg=userMessage{}
				result int
				err error
				user model.Users
				str []byte
				ctx = websocket.GetContext(nsConn.Conn)
			)

			err=msg.Unmarshal(&userMsg)
			if err!=nil{
				return err
			}

			user=ctx.Values().Get("user").(model.Users)

			result,err=service.NewUsersUsersService().EachOtherFriends(user.AppsId,userMsg.To,userMsg.From)
			if result!=service.EACH_OTHER_FRIENDS_IS_TRUE{
				//校验失败时返回消息
				tempErrCode:=common.CONTACTS_BROKEN
				userMsg.ErrCode=&tempErrCode
				tempErrMsg:="消息发送失败,您不在对方好友列表."
				userMsg.ErrMsg=&tempErrMsg

				str,_=userMsg.Marshal()
				nsConn.Emit("chatTo",str)
				return nil
			}

			msg.To= strconv.FormatUint(userMsg.To,10)
			log.Println(string(msg.Body))
			nsConn.Conn.Server().Broadcast(nil,msg)
			return nil
		},

	}
}