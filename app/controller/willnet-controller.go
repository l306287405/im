package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"log"
)

func WillnetController() websocket.Events{
	return websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)
			jwtStr := ctx.Values().Get("jwt")
			if jwtStr == nil{
				nsConn.Conn.Close()
				return errors.New("用户令牌失效,关闭连接")
			}
			user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)


			fmt.Println("登录用户信息:",user)

			apps_id,user_id:=uint(user["apps_id"].(float64)),uint64(user["id"].(float64))
			err := dao.NewUsersDao().Online(apps_id,user_id,nsConn.Conn.ID())
			if err!=nil{
				r,_:=json.Marshal(common.SendCry("用户上线失败,关闭链接"))
				nsConn.Emit("chat",r)
				nsConn.Conn.Close()
				return errors.New("用户上线失败,关闭链接")
			}

			//room_id:=int(user["id"].(float64))%2
			//room:=nsConn.Room(strconv.Itoa(room_id))
			room,err:=nsConn.JoinRoom(nil,"test")
			if err!=nil{
				 nsConn.Conn.Close()
				 return errors.New("加入房间失败,关闭连接")
			}

			room.Emit("chat",[]byte(fmt.Sprintf("欢迎 [%s] 加入 [%s] 号房间",user["nickname"].(string),room.Name)))

			//nsConn.Conn.Close()

			log.Printf("ConnID [%s] connected to namespace [%s] with IP [%s],Nickname [%s],ID [%f]",
				nsConn, msg.Namespace, ctx.RemoteAddr(),user["nickname"],user["id"])
			return nil
		},

		//断开连接
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {

			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},

		//websocket.OnRoomJoin: func(nsConn *neffos.NSConn, message neffos.Message) error {
		//	log.Printf("server: 接入房间 %s", message.Room)
		//	return nil
		//},

		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			jwtStr := ctx.Values().Get("jwt")
			if jwtStr == nil{
				return errors.New("用户令牌失效,关闭连接")
			}

			user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)
			//room_id:=int(user["id"].(float64))%2
			//room,err:=nsConn.JoinRoom(nil,strconv.Itoa(room_id))

			room:=nsConn.Room("test")

			//oom.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			msg_body:=fmt.Sprintf("[%s] in [%s,%s] sent: %s", user["nickname"],room.Name,room, string(msg.Body))
			log.Println(msg_body)

			// Write message back to the client message owner with:
			//nsConn.Emit("chat", msg.Body)
			// Write message to all except this client with:
			room.Emit("chat",[]byte(msg_body))
			//nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},

	}
}