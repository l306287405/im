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

	//app授权
	mvc.Configure(app.Party("/apps"), func(app *mvc.Application) {
		app.Handle(new(controller.AppsController))
	})

	//用户权限与注册相关接口
	mvc.Configure(app.Party("/users"), func(app *mvc.Application) {
		//中间件
		app.Router.Use(Middleware.AppsCheck)

		//依赖注入
		app.Register(
			sessions.New(sessions.Config{}).Start,
		)

		//登录注册
		app.Handle(new(controller.UsersAccessController))

	})

	//联系人相关接口
	mvc.Configure(app.Party("/users/contacts/users/{fid:uint64}"), func(app *mvc.Application) {
		app.Router.Use(Middleware.UsersVerify,Middleware.AppsCheck)

		app.Handle(new(controller.UsersContactsUsersController))
	})

	//群聊相关接口
	mvc.Configure(app.Party("/chatrooms"), func(app *mvc.Application) {
		app.Router.Use(Middleware.UsersVerify,Middleware.AppsCheck)
		app.Handle(new(controller.ChatroomsController))
	})
	mvc.Configure(app.Party("/chatrooms/{room_id:uint64}/users"), func(app *mvc.Application) {
		app.Router.Use(Middleware.UsersVerify,Middleware.AppsCheck)
		app.Handle(new(controller.ChatroomsUsersController))
	})

	//
	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString(time.Now().String())
	})

}
