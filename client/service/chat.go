package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"im/client/common"
	"im/client/model"
	"log"
	"os"
	"strconv"
	"time"
)

type Chat struct {
	u *model.Users

}

type Rooms map[string]*websocket.Room
var rooms Rooms

func NewChat() *Chat{
	return &Chat{}
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
	websocket.OnRoomJoined: func(nsConn *websocket.NSConn, msg websocket.Message) error {
		log.Printf("%s connected to room %s", nsConn,msg.Room)
		return nil
	},
	websocket.OnRoomLeft: func(nsConn *websocket.NSConn, msg websocket.Message) error {
		log.Printf("%s left from room %s", nsConn,msg.Room)
		return nil
	},
	"group": func(c *websocket.NSConn, msg websocket.Message) error {
		log.Println(string(msg.Serialize()))
		//var userMsg model.GroupsMessages
		//err := msg.Unmarshal(&userMsg)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println(userMsg)
		return nil
	},
	"chat": func(c *websocket.NSConn, msg websocket.Message) error {
		fmt.Println(string(msg.Serialize()))
		return nil
	},
	"notify": func(c *websocket.NSConn, msg websocket.Message) error {
		fmt.Println(string(msg.Serialize()))
		return nil
	},
	"receipt": func(nsConn *websocket.NSConn, msg websocket.Message) error {
		fmt.Println(string(msg.Serialize()))
		return nil
	},

}

func (ch *Chat) Connect(appToken string,authToken string,appsId uint,userId uint64){
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(common.DialAndConnectTimeout))
	defer cancel()

	events := make(websocket.Namespaces)
	events["W"]=clientEvents

	// init the websocket connection by dialing the server.
	client, err := websocket.Dial(ctx, websocket.GobwasDialer(websocket.GobwasDialerOptions{}),

		fmt.Sprintf("%s?X-Websocket-Header-X-APP-Token=%s&X-Websocket-Header-X-JWT-Token=%s",common.Endpoint,appToken,authToken),
		// The namespaces and events, can be optionally shared with the server's.
		events)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	c, err := client.Connect(ctx, "W")
	if err != nil {
		panic(err)
	}

	//Get list of groups and join they
	result,err:=common.HttpDo("GET",fmt.Sprintf("%s/chatrooms?X-Websocket-Header-X-JWT-Token=%s",
		common.Host,authToken),nil,nil)
	if err!=nil{
		fmt.Fprintln(os.Stdout,"获取群列表失败:"+err.Error())
		return
	}
	var rep struct{
		Msg string `json:"msg"`
		Code int	`json:"code"`
		Data []model.ChatroomsUsers `json:"data"`
	}
	err=json.Unmarshal(result,&rep)
	if err!=nil{
		fmt.Fprintln(os.Stdout,"获取群列表结果不符合规范",rep,result)
		return
	}

	var roomsListStr []uint64
	if len(rep.Data)>0{
		rooms=make(Rooms)
	}
	for _,v:=range rep.Data{
		room, err := c.JoinRoom(nil, strconv.FormatUint(v.RoomId,10))
		if err != nil {
			log.Fatal(err)
		}
		rooms[strconv.FormatUint(v.RoomId,10)]=room
		roomsListStr=append(roomsListStr, v.RoomId)
	}

	var chatType string
	var b []byte
BEGIN:
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintln(os.Stdout,"私聊输入1,群聊输入2")
	if !scanner.Scan(){
		log.Printf("ERROR: %v", scanner.Err())
		return
	}
	chatType = scanner.Text()

	var obj string
	var roomNow *websocket.Room

RE:

	if chatType=="1" {
		fmt.Fprintln(os.Stdout,"输入你想要私聊的用户")
		if !scanner.Scan(){
			log.Printf("ERROR: %v", scanner.Err())
			return
		}
		obj = scanner.Text()
	}else {
		fmt.Fprint(os.Stdout,"当前您已加入的房间:")
		fmt.Fprintln(os.Stdout,roomsListStr)
		fmt.Fprintln(os.Stdout,"输入你要加入群聊的房间号.")
		if !scanner.Scan(){
			log.Printf("ERROR: %v", scanner.Err())
			return
		}
		obj = scanner.Text()
		if _,ok:=rooms[obj];ok{
			roomNow=rooms[obj]
		}else{
			fmt.Fprintln(os.Stdout,"您输入的房间号不在您的房间列表中")
			goto RE
		}
	}

	fmt.Fprint(os.Stdout, ">> ")
	for {
		if !scanner.Scan() {
			log.Printf("ERROR: %v", scanner.Err())
			return
		}

		text := scanner.Text()

		var textType = model.MSG_TYPE_IS_TEXT
		switch text {
		case "exit":
			goto BEGIN

		case "leave":
			roomNow.Leave(nil)

		case "mock_img":
			m:=model.Multimedia{
				Id:   2,
				Url:  "http://wj-app.oss-cn-hangzhou.aliyuncs.com/uploads/image/jpeg/20190716/13c95def7935730d12b3adb0f2b73ef8f06e51bd.jpeg?OSSAccessKeyId=LTAITEDzaZ1V6Zz7&Expires=1601072253&Signature=tIqPL31xfVGsX%2BK5EGZEFnWfEnw%3D",
				Mime: "image/jpeg",
				Size: 73847,
			}
			b,_=json.Marshal(m)
			text=string(b)
			textType=model.MSG_TYPE_IS_IMAGE
		case "mock_video":
			m:=model.Multimedia{
				Id:   3,
				Url:  "https://wj-app.oss-cn-hangzhou.aliyuncs.com/uploads/video/mp4/20190717/9d29759cb743c69df26d536238e8045664ea6dc3.mp4?Expires=1717118572&OSSAccessKeyId=LTAIWcz38433BCtR&Signature=v5YXrniEqzfi9I/pitarATrRYGo%3D&security-token",
				Mime: "video/mp4",
				Size: 3555256,
			}
			b,_=json.Marshal(m)
			text=string(b)
			textType=model.MSG_TYPE_IS_VIDEO
		}


		to,_:=strconv.ParseUint(obj,10,64)
		userMsg := model.Messages{AppsId:appsId,From: userId,To:to, Text: text,TextType:textType,Status:1}
		if chatType=="1"{
			c.Emit("chat",websocket.Marshal(userMsg))
		}else{
			roomNow.Emit("group", websocket.Marshal(userMsg))
		}


		fmt.Fprint(os.Stdout, ">> ")
	}

}