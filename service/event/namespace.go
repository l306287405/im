package event

import (
	"encoding/json"
	"errors"
	"github.com/kataras/iris/websocket"
	"im/common"
	"im/dao"
	"im/model"
	"log"
)

func OnNamespaceConnected() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
		ctx := websocket.GetContext(nsConn.Conn)
		user:=ctx.Values().Get("user").(model.Users)

		//加入在线人数
		err := dao.NewUsersDao().Online(user.AppsId,user.Id,nsConn.Conn.ID())
		if err!=nil{
			r,_:=json.Marshal(common.SendCry("用户上线失败,关闭链接"))
			nsConn.Emit("chat",r)
			nsConn.Conn.Close()
			return errors.New( "用户上线失败,关闭链接")
		}

		// with `websocket.GetContext` you can retrieve the Iris' `Context`.
		log.Printf("[%s] connected to namespace [%s].",nsConn, msg.Namespace)
		return nil
	}
}

func OnNamespaceDisconnect() func(*websocket.NSConn,websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
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
			nsConn.Emit("chat",r)
			nsConn.Conn.Close()
			return errors.New("用户下线失败,关闭链接")
		}

		log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
		return nil
	}
}

