package router

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"im/app/controller"
	"log"
	"net/http"
	"os"
)

// if namespace is empty then simply websocket.Events{...} can be used instead.

func WebsocketRouter(app *iris.Application) {

	//构建服务
	events := make(websocket.Namespaces)
	events["willnet"]=controller.WillnetController()

	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, events)

	//id命名
	websocketServer.IDGenerator = func(w http.ResponseWriter, r *http.Request) string {

		jwtStr:=r.Header.Get("jwt")
		token:=new(jwt.Token)
		if jwtStr == ""{
			return websocket.DefaultIDGenerator(nil)
		}
		err:=json.Unmarshal([]byte(r.Header.Get("jwt")),token)
		if err!=nil{
			return websocket.DefaultIDGenerator(nil)
		}

		user := token.Claims.(jwt.MapClaims)
		return user["nickname"].(string)
	}

	//连接升级失败
	websocketServer.OnUpgradeError = func(err error) {
		log.Printf("ERROR: %v", err)
	}

	//连接状态
	websocketServer.OnConnect = func(c *neffos.Conn) error {
		if c.WasReconnected() {
			log.Printf("[%s] connection is a result of a client-side re-connection, with tries: %d", c.ID(), c.ReconnectTries)
		}

		log.Printf("[%s] connected to the server.", c)

		// if returns non-nil error then it refuses the client to connect to the server.
		return nil
	}

	websocketServer.OnDisconnect = func(c *neffos.Conn) {
		log.Printf("[%s] disconnected from the server.", c)
	}

	fmt.Println(os.Getenv("JWT_SECRET"))
	websocketRouter := app.Get("/websocket/echo", websocket.Handler(websocketServer))

	//jwt
	jwtHandler := jwt.New(jwt.Config{
		Extractor:jwt.FromAuthHeader,
		ValidationKeyGetter: func(token *jwt.Token) (i interface{}, e error) {
			return []byte(os.Getenv("JWT_SECRET")),nil
		},
		SigningMethod:jwt.SigningMethodHS256,
	})


	websocketRouter.Use(jwtHandler.Serve)
}
