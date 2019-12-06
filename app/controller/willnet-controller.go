package controller

import (
	"github.com/kataras/iris/websocket"
	"im/service/event"
)

func W() websocket.Events{
	return websocket.Events{

		//连入命名空间
		websocket.OnNamespaceConnect: event.OnNamespaceConnect(),
		websocket.OnNamespaceConnected: event.OnNamespaceConnected(),
		websocket.OnNamespaceDisconnect: event.OnNamespaceDisconnect(),

		//连入房间
		websocket.OnRoomJoin: event.OnRoomJoin(),
		websocket.OnRoomJoined: event.OnRoomJoined(),
		websocket.OnRoomLeft: event.OnRoomLeft(),

		//群聊
		"group": event.Group(),

		//私聊
		"chat": event.Chat(),
	}
}