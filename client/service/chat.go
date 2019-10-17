package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/neffos"
	"im/client/common"
	"im/model"
	"log"
	"os"
	"time"
)

type Chat struct {
	u *model.Users

}

func NewChat() *Chat{
	return &Chat{}
}
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


//客户端的节点回调函数
var clientEvents = websocket.Events{
	websocket.OnNamespaceConnected: func(c *websocket.NSConn, msg websocket.Message) error {
		log.Printf("connected to namespace: %s", msg.Namespace)
		return nil
	},
	websocket.OnNamespaceDisconnect: func(c *websocket.NSConn, msg websocket.Message) error {
		log.Printf("disconnected from namespace: %s", msg.Namespace)
		return nil
	},
	websocket.OnRoomJoined: func(nsConn *neffos.NSConn, msg neffos.Message) error {
		log.Printf("%s 接入房间 %s", nsConn,msg.Room)
		return nil
	},
	websocket.OnRoomLeft: func(nsConn *neffos.NSConn, msg neffos.Message) error {
		log.Printf("%s 离开房间 %s", nsConn,msg.Room)
		return nil
	},
	"chat": func(c *websocket.NSConn, msg websocket.Message) error {
		var userMsg userMessage
		err := msg.Unmarshal(&userMsg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s >> [%s] says: %s\n", msg.Room, userMsg.From, userMsg.Text)
		return nil
	},
	"chatTo": func(c *websocket.NSConn, msg websocket.Message) error {
		fmt.Println(string(msg.Body))
		return nil
	},

}

func (chat *Chat) Connect(authToken string,userId uint64){
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(common.DialAndConnectTimeout))
	defer cancel()

	events := make(websocket.Namespaces)
	events["willnet"]=clientEvents

	// init the websocket connection by dialing the server.
	client, err := websocket.Dial(
		// Optional context cancelation and deadline for dialing.
		ctx,
		// The underline dialer, can be also a gobwas.Dialer/DefautlDialer or a gorilla.Dialer/DefaultDialer.
		// Here we wrap a custom gobwas dialer in order to send the username among, on the handshake state,
		// see `startServer().server.IDGenerator`.
		websocket.GobwasDialer(websocket.GobwasDialerOptions{Header:websocket.GobwasHeader{"Authorization": []string{"Bearer "+authToken}}}),
		// The endpoint, i.e ws://localhost:8080/path.
		common.Endpoint,
		// The namespaces and events, can be optionally shared with the server's.
		events)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	c, err := client.Connect(ctx, "willnet")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintln(os.Stdout,"私聊输入1,群聊输入2")
	if !scanner.Scan(){
		log.Printf("ERROR: %v", scanner.Err())
		return
	}
	chatType := scanner.Text()

	obj:=""
	var room=&websocket.Room{}
	if chatType == "1" {
		fmt.Fprintln(os.Stdout,"输入你想要私聊的用户")
		if !scanner.Scan(){
			log.Printf("ERROR: %v", scanner.Err())
			return
		}
		obj = scanner.Text()
	}else {
		fmt.Fprintln(os.Stdout,"输入你要加入群聊的房间")
		if !scanner.Scan(){
			log.Printf("ERROR: %v", scanner.Err())
			return
		}
		obj = scanner.Text()
		room, err = c.JoinRoom(nil, obj)
		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Fprint(os.Stdout, ">> ")
	for {
		if !scanner.Scan() {
			log.Printf("ERROR: %v", scanner.Err())
			return
		}

		text := scanner.Text()

		if text=="exit" {
			if err := c.Disconnect(nil); err != nil {
				log.Printf("reply from server: %v", err)
			}
			break
		}

		if text == "leave" {
			room.Leave(nil)
			break
		}


		userMsg := userMessage{From: fmt.Sprintf("%d",userId),To:obj,Type:chatType, Text: text}
		if chatType == "1"{

			c.Emit("chatTo",websocket.Marshal(userMsg))
		}else{
			room.Emit("chat", websocket.Marshal(userMsg))
		}


		fmt.Fprint(os.Stdout, ">> ")
	}

}