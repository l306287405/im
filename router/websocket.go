package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
	"im/app/controller"
	"im/model"
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

	//启动服务,定义id生成规则,将用户信息存储至连接
	websocketRouter := app.Get("/websocket/echo", websocket.Handler(websocketServer, func(ctx context.Context) string {
		jwtStr:=ctx.Values().Get("jwt")
		log.Println(jwtStr)
		if jwtStr == nil{
			return websocket.DefaultIDGenerator(ctx)
		}

		user := jwtStr.(*jwt.Token).Claims.(jwt.MapClaims)

		loginUser:=new(model.Users)
		loginUser.Id=uint64(user["id"].(float64))
		loginUser.AppsId=uint(user["apps_id"].(float64))
		loginUser.Account=user["account"].(string)
		loginUser.Nickname=user["nickname"].(string)
		ctx.Values().Set("user",*loginUser)

		return strconv.FormatUint(uint64(user["id"].(float64)),10)
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
