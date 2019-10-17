package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
	"im/app/controller"
	"log"
	"os"
	"strconv"
)

// if namespace is empty then simply websocket.Events{...} can be used instead.

func WebsocketRouter(app *iris.Application) {

	//构建服务
	events := make(websocket.Namespaces)
	events["willnet"]=controller.Willnet()

	websocketServer := websocket.New(websocket.DefaultGobwasUpgrader, events)

	//连接升级失败
	websocketServer.OnUpgradeError = func(err error) {
		log.Printf("ERROR: %v", err)
	}

	//连接状态
	websocketServer.OnConnect = func(c *websocket.Conn) error {
		if c.WasReconnected() {
			log.Printf("[%s] connection is a result of a client-side re-connection, with tries: %d", c.ID(), c.ReconnectTries)
		}

		log.Printf("[%s] connected to the server.", c)

		// if returns non-nil error then it refuses the client to connect to the server.
		return nil
	}

	websocketServer.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] disconnected from the server.", c)
	}

	//启动服务并定义id生成规则
	websocketRouter := app.Get("/websocket/echo", websocket.Handler(websocketServer, func(ctx context.Context) string {
		jwtStr:=ctx.Values().Get("jwt")
		log.Println(jwtStr)
		if jwtStr == nil{
			return websocket.DefaultIDGenerator(ctx)
		}

		user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)

		return strconv.FormatFloat(user["apps_id"].(float64),'f',-1,64)+"_"+user["account"].(string)
	}))

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
