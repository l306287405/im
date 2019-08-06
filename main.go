package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"im/router"
	"im/service/orm"
	"log"
	"time"
)

func main(){
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	//优雅关闭
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx)
	})
	app.ConfigureHost(func(h *iris.Supervisor) {
		h.RegisterOnShutdown(func() {
			println("服务已关闭")
		})
	})

	//环境配置
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//数据库结构同步
	orm.SyncDB()

	//静态资源
	app.HandleDir("/static","./static")

	//构建路由
	router.Run(app)

	//启动应用
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./configs/iris.tml")))

}