package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"im/model"
	"im/service"
	"im/service/orm"
	"log"
	"strconv"
	"sync"
)

type roomsUsersTable struct {
	mutex	 sync.RWMutex
	entries  map[uint64][]*websocket.NSConn
}

func (ru *roomsUsersTable) GetLen(id uint64) int{
	l := 0
	ru.mutex.RLock()
	if _,ok:=ru.entries[id];ok{
		l=len(ru.entries[id])
	}
	ru.mutex.RUnlock()
	return l
}

func (ru *roomsUsersTable) AddConnOfRoom(id uint64,conn *websocket.NSConn) bool{

	if !ru.RoomExist(id){
		return false
	}
	if ru.UserExist(id,conn){
		return false
	}
	ru.mutex.Lock()
	ru.entries[id]=append(ru.entries[id], conn)
	ru.mutex.Unlock()
	return true
}

func (ru *roomsUsersTable) DelConnOfRoom(id uint64,conn *websocket.NSConn) bool{
	if !ru.RoomExist(id){
		return false
	}
	ru.mutex.Lock()
	defer ru.mutex.Unlock()
	for i, v := range ru.entries[id] {
		if v == conn {
			ru.entries[id] = append(ru.entries[id][:i], ru.entries[id][i+1:]...)
			return true
		}
	}
	return false

}

func (ru *roomsUsersTable) RoomExist(id uint64) bool{
	ru.mutex.RLock()
	_,ok:=ru.entries[id]
	ru.mutex.RUnlock()
	return ok
}

func (ru *roomsUsersTable) UserExist(id uint64,conn *websocket.NSConn) bool{
	if !ru.RoomExist(id){
		return false
	}
	ru.mutex.RLock()
	defer ru.mutex.RUnlock()
	for _,v:=range ru.entries[id]{
		if v==conn{
			return true
		}
	}

	return false
}

var RoomUsersTable roomsUsersTable

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
				return errors.New( "用户上线失败,关闭链接")
			}

			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			log.Printf("[%s] connected to namespace [%s].",nsConn, msg.Namespace)
			return nil
		},

		//断开连接
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			if ctx.Values().Get("user")==nil{
				log.Printf("用户信息不正确")
				return nil
			}
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

		websocket.OnRoomJoined: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			var(
				ctx = websocket.GetContext(nsConn.Conn)
				user=ctx.Values().Get("user").(model.Users)
				userMsg=model.Messages{}
				err=msg.Unmarshal(&userMsg)
				roomId uint64
				status *int8
				tempCode string
				tempMsg string
				str []byte
			)

			roomId,err=strconv.ParseUint(msg.Room,10,64)
			if err!=nil{
				return err
			}
			status=dao.NewChatroomsUsersDao().RelationExist(user.AppsId,roomId,user.Id)
			if status==nil || *status!=1{
				tempCode=common.ROOMS_USERS_BROKEN
				userMsg.ErrCode=&tempCode
				tempMsg="非群聊成员无法加入"
				userMsg.ErrMsg=&tempMsg
				str,_=userMsg.Marshal()
				nsConn.Emit("chatTo",str)
				err=nsConn.Room(msg.Room).Leave(nil)
				if err!=nil{
					return err
				}
				return nil
			}

			//加入房间
			RoomUsersTable.AddConnOfRoom(roomId,nsConn)

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

			roomId,err:=strconv.ParseUint(msg.Room,10,64)
			if err!=nil{
				return err
			}
			RoomUsersTable.DelConnOfRoom(roomId,nsConn)

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
			var(
				userMsg=model.GroupsMessages{}
				ctx = websocket.GetContext(nsConn.Conn)
				userInter=ctx.Values().Get("user")
				user=model.Users{}
				str []byte
				err error
				tempCode string
				tempMsg string
				db=orm.GetDB()
			)
			if userInter==nil{
				fmt.Println("用户信息有误,请检查")
				return nil
			}
			user=userInter.(model.Users)

			fmt.Println(string(msg.Serialize()))
			err=msg.Unmarshal(&userMsg)
			if err!=nil{
				return err
			}

			//插入消息记录
			userMsg.AppsId,userMsg.Status=user.AppsId,1
			_,err=db.InsertOne(&userMsg)
			if err!=nil{
				tempCode=common.SQL_INSERT_FAILD
				tempMsg=err.Error()
				userMsg.ErrCode,userMsg.ErrMsg=&tempCode,&tempMsg
				str,_=userMsg.Marshal()
				nsConn.Emit("chat",str)
				return nil
			}

			//TODO act报文反馈 已发送 已阅读
			str,_=userMsg.Marshal()
			msg.Body=str
			nsConn.Conn.Server().Broadcast(nsConn, msg)

			return nil
		},
		"chatTo": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			var(
				userMsg=model.Messages{}
				result int
				err error
				ctx = websocket.GetContext(nsConn.Conn)
				user =ctx.Values().Get("user").(model.Users)
				str []byte
				db=orm.GetDB()
				tempCode string
				tempMsg string
			)

			fmt.Println(string(msg.Serialize()))
			err=msg.Unmarshal(&userMsg)
			if err!=nil{
				return err
			}

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

			userMsg.AppsId,userMsg.Status=user.AppsId,1
			_,err=db.InsertOne(&userMsg)
			if err!=nil{
				tempCode=common.SQL_INSERT_FAILD
				tempMsg=err.Error()
				userMsg.ErrCode,userMsg.ErrMsg=&tempCode,&tempMsg
				str,_=userMsg.Marshal()
				nsConn.Emit("chat",str)
				return nil
			}

			//TODO act报文反馈 已发送 已阅读
			str,_=userMsg.Marshal()
			msg.Body=str

			nsConn.Conn.Server().Broadcast(nil,msg)
			return nil
		},

	}
}