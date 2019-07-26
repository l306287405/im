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
		//当然，你可以在MVC应用程序中使用普通的中间件。
		app.Router.Use(Middleware.AppsCheck)
		//把依赖注入，controller(s)绑定
		//可以是一个接受iris.Context并返回单个值的函数（动态绑定）
		//或静态结构值（service）。
		app.Register(
			sessions.New(sessions.Config{}).Start,
		)
		// GET: http://localhost:8080/basic
		// GET: http://localhost:8080/basic/custom
		app.Handle(new(controller.UserAccessController))
		//所有依赖项被绑定在父 *mvc.Application
		//被克隆到这个新子身上，父的也可以访问同一个会话。
		// GET: http://localhost:8080/basic/sub
		//app.Party("/sub").Handle(new(basicSubController))

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
