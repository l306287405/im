package router

import (
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"im/app/controller"
	"os"
)

// if namespace is empty then simply websocket.Events{...} can be used instead.

func WebsocketRouter(app *iris.Application) {

	//构建服务
	events := make(websocket.Namespaces)
	events["willnet"]=controller.WillnetController()

	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, events)


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

	//websocketServer.OnConnect = func(c *websocket.Conn) error {
	//	ctx := websocket.GetContext(c)
	//	if err := jwtHandler.CheckJWT(ctx); err != nil {
	//		// will send the above error on the client
	//		// and will not allow it to connect to the websocket server at all.
	//		return err
	//	}
	//	user := ctx.Values().Get("jwt").(*jwt.Token)
	//	fmt.Println(user)
	//	// or just: user := j.Get(ctx)
	//	log.Printf("This is an authenticated request\n")
	//	log.Printf("Claim content:")
	//	log.Printf("%#+v\n", user.Claims)
	//	log.Printf("[%s] connected to the server", c.ID())
	//	return nil
	//}


	websocketRouter.Use(jwtHandler.Serve)
}
