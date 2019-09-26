package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/neffos/gorilla"
	"im/client/common"
	"im/model"
	"log"
	"net/http"
	"os"
	"time"
)

type Chat struct {
	u *model.Users

}

func NewChat() *Chat{
	return &Chat{}
}

// this can be shared with the server.go's.
// `NSConn.Conn` has the `IsClient() bool` method which can be used to
// check if that's is a client or a server-side callback.
var clientEvents = websocket.Namespaces{
	"willnet": websocket.Events{
		websocket.OnNamespaceConnected: func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("connected to namespace: %s", msg.Namespace)
			return nil
		},
		websocket.OnNamespaceDisconnect: func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("disconnected from namespace: %s", msg.Namespace)
			return nil
		},
		//websocket.OnRoomJoined: func(nsConn *neffos.NSConn, msg neffos.Message) error {
		//	log.Printf("%s 接入房间 %s", nsConn,msg.Room)
		//	return nil
		//},
		//websocket.OnRoomLeft: func(nsConn *neffos.NSConn, msg neffos.Message) error {
		//	log.Printf("%s 离开房间 %s", nsConn,msg.Room)
		//	return nil
		//},
		"chat": func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("%s", string(msg.Body))
			return nil
		},
	},
}

func (chat *Chat) Connect(authToken string,userId uint64){
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(common.DialAndConnectTimeout))
	defer cancel()

	client, err := websocket.Dial(ctx, gorilla.Dialer(&gorilla.Options{},http.Header{"Authorization": []string{"Bearer "+authToken}}), common.Endpoint, clientEvents)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	c, err := client.Connect(ctx, "willnet")
	if err != nil {
		panic(err)
	}

	c.Emit("chat", []byte("Hello from Go client side!"))
	//c.Room("1").NSConn.Conn.

	fmt.Fprint(os.Stdout, ">> ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			log.Printf("ERROR: %v", scanner.Err())
			return
		}

		text := scanner.Bytes()

		if bytes.Equal(text, []byte("exit")) {
			if err := c.Disconnect(nil); err != nil {
				log.Printf("reply from server: %v", err)
			}
			break
		}

		fmt.Println(c.Rooms())
		//ok := c.Emit("chat", text)
		//room_id:=int(userId)%2

		room:=c.Room("test")

		ok:=room.NSConn.Emit("chat",text)

		if !ok {
			break
		}

		fmt.Fprint(os.Stdout, ">> ")
	}

}