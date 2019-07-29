package router

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"im/app/Middleware"
	"im/app/controller"
	"time"
)

func AppRouter(app *iris.Application){

	mvc.Configure(app.Party("/users"), func(app *mvc.Application) {
		//中间件
		app.Router.Use(Middleware.AppsCheck)

		//依赖注入
		app.Register(
			sessions.New(sessions.Config{}).Start,
		)

		app.Handle(new(controller.UserAccessController))

	})

	mvc.Configure(app.Party("/apps"), func(app *mvc.Application) {
		app.Handle(new(controller.AppsController))
	})


	//Method:   GET
	//Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString(time.Now().String())
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})


	//app.Post("/login",)

	////分组路由示例
	//app.PartyFunc("/users",func(users iris.Party){
	//	users.Use(myAuthMiddlewareHandler)
	//	users.Get("/inbox/{id:int}",handler)
	//})

}
