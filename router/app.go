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

	//注册视图
	app.RegisterView(iris.HTML("./app/view",".html"))

	//用户登录注册
	mvc.Configure(app.Party("/users"), func(app *mvc.Application) {
		//中间件
		app.Router.Use(Middleware.AppsCheck)

		//依赖注入
		app.Register(
			sessions.New(sessions.Config{}).Start,
		)

		app.Handle(new(controller.UserAccessController))

	})

	//app授权
	mvc.Configure(app.Party("/apps"), func(app *mvc.Application) {
		app.Handle(new(controller.AppsController))
	})


	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString(time.Now().String())
	})

}
