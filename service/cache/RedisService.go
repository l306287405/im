package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"sync"
)

const(

	//app与token的缓存 hash映射	格式:iota_keyId_keySecret => token
	AppsTokensMap = iota

)

var (
	instance *redis.Client
	once     sync.Once
)

func Get() *redis.Client {
	once.Do(func() {
		var (
			host     = os.Getenv("REDIS_HOST")
			port     = os.Getenv("REDIS_PORT")
			password = os.Getenv("REDIS_PASSWORD")
			db, _    = strconv.Atoi(os.Getenv("REDIS_DB"))
			address  = fmt.Sprintf("%s:%s", host, port)
		)

		instance = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		})

		pong,err := instance.Ping().Result()

		if err != nil {
			panic("Redis连接失败 原因:"+err.Error())
		}

		println("Redis启动:"+pong)
	})
	return instance
}
