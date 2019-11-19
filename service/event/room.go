package event

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"im/model"
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

func OnRoomJoined() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
		var(
			ctx = websocket.GetContext(nsConn.Conn)
			user=ctx.Values().Get("user").(model.Users)
			userMsg=model.GroupsMessages{}
			err=msg.Unmarshal(&userMsg)
			roomId uint64
			status *int8
		)

		roomId,err=strconv.ParseUint(msg.Room,10,64)
		if err!=nil{
			return err
		}
		status=dao.NewChatroomsUsersDao().RelationExist(user.AppsId,roomId,user.Id)
		if status==nil || *status!=1{
			//tempCode=common.ATTRIBUTION_ERROR
			//userMsg.ErrCode=&tempCode
			//tempMsg="非群聊成员无法加入"
			//userMsg.ErrMsg=&tempMsg
			userMsg.Err=errors.New(common.ATTRIBUTION_ERROR)
			//str,_=userMsg.Marshal()
			nsConn.Emit("chat",websocket.Marshal(userMsg))
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
	}
}

func OnRoomLeft() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
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
	}
}
