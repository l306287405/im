package controller

import (
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/websocket"
	"log"
)

func WillnetController() websocket.Events{
	return websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			user := ctx.Values().Get("jwt").(*jwt.Token)
			fmt.Println(user)
			//nsConn.Conn.Close()

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))

			// Write message back to the client message owner with:
			//nsConn.Emit("chat", msg.Body)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},

	}
}