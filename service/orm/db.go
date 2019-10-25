package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"im/model"
	"os"
	"sync"
)

var (
	instance *xorm.Engine
	once sync.Once
)

func GetDB() *xorm.Engine {

	once.Do(func() {
		host := os.Getenv("DB_HOST")
		port :=os.Getenv("DB_PORT")
		database :=os.Getenv("DB_DATABASE")
		username :=os.Getenv("DB_USERNAME")
		password :=os.Getenv("DB_PASSWORD")
		charset := os.Getenv("DB_CHARSET")
		db_config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",username,password,host,port,database,charset)
		var err error
		db, err := xorm.NewEngine("mysql",db_config)
		if err != nil {
			panic("数据库连接失败")
		}

		instance = db
	})
	return instance
}

//数据库结构同步方法
func SyncDB(){
		db:=GetDB()
		err := db.Sync2(
			new(model.Apps),
			new(model.Chatrooms),
			new(model.ChatroomsUsers),
			new(model.Classes),
			new(model.Files),
			new(model.Messages),
			new(model.Users),
			new(model.UsersUsers))
	if err!=nil{
		panic("数据库结构同步失败 原因:"+err.Error())
	}
}