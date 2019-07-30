package router

import "github.com/kataras/iris"

func Run(app *iris.Application){

	//注册常规api
	AppRouter(app)

	//注册长连接
	WebsocketRouter(app)
}