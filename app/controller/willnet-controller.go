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

// userMessage implements the `MessageBodyUnmarshaler` and `MessageBodyMarshaler`.
type userMessage struct {
	From string `json:"from"`
	To	 string	`json:"to"`
	Type string	`json:"type"`
	Text string `json:"text"`
}

// Defaults to `DefaultUnmarshaler & DefaultMarshaler` that are calling the json.Unmarshal & json.Marshal respectfully
// if the instance's Marshal and Unmarshal methods are missing.
func (u *userMessage) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *userMessage) Unmarshal(b []byte) error {
	return json.Unmarshal(b, u)
}

func Willnet() websocket.Events{
	return websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {

			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			log.Printf("[%s] connected to namespace [%s].",nsConn, msg.Namespace)
			return nil
		},

		//断开连接
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {

			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},

		websocket.OnRoomJoined: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			ctx := websocket.GetContext(nsConn.Conn)
			jwtStr := ctx.Values().Get("jwt")
			if jwtStr == nil{
				nsConn.Conn.Close()
				return errors.New("用户令牌失效,关闭连接")
			}
			user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)

			apps_id,user_id:=uint(user["apps_id"].(float64)),uint64(user["id"].(float64))
			err := dao.NewUsersDao().Online(apps_id,user_id,nsConn.Conn.ID())
			if err!=nil{
				r,_:=json.Marshal(common.SendCry("用户上线失败,关闭链接"))
				nsConn.Emit("chat",r)
				nsConn.Conn.Close()
				return errors.New("用户上线失败,关闭链接")
			}

			text := fmt.Sprintf("欢迎ID [%d] 的用户 [%s] 加入 [%s] 号房间",user_id,user["nickname"], msg.Room)
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
			ctx := websocket.GetContext(nsConn.Conn)
			jwtStr := ctx.Values().Get("jwt")
			if jwtStr == nil{
				return errors.New("用户令牌失效,关闭连接")
			}

			//user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)

			nsConn.Conn.Server().Broadcast(nsConn, msg)

			//room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			//msg_body:=fmt.Sprintf("[%s] in [%s,%s] sent: %s", user["nickname"],room.Name,room, string(msg.Body))
			//log.Println(msg_body)

			// Write message back to the client message owner with:
			//nsConn.Emit("chat", msg.Body)
			// Write message to all except this client with:
			//nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
		"chatTo": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			var(
				userMsg=userMessage{}
			)

			err:=msg.Unmarshal(&userMsg)
			if err!=nil{
				return err
			}
			msg.To=userMsg.To
			msg.FromExplicit=nsConn.Conn.ID()
			nsConn.Conn.Server().Broadcast(nil,msg)
			return nil
		},

	}
}