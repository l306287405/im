package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"sync"
)

const(

	//app与token的缓存	格式:iota_token => appId
	APPS_TOKEN_MAP = iota

	//用户token缓存	格式:iota_appId_userToken => userId
	USERS_TOKEN_MAP

	//用户与通讯id的映射缓存 格式:iota_appId_userId => cId
	USERS_COMM_MAP

	//用户与用户互为好友的关系映射 格式:iota_appId_userId(小值)_userId(大值) => 1
	EACH_OTHER_FRIENDS

	//群聊房间与用户的关系映射 格式:iota_appId_roomId_userId => relationStatus -1:软删 0:待审核 1:正常
	ROOMS_USERS_MAP

)

var (
	instance *redis.Client
	once     sync.Once
)

func Init() *redis.Client {
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
			return
		}

		println("Redis启动:"+pong)
	})
	return instance
}

//统一
func CacheKey(i int,s string) string{
	return fmt.Sprintf("%d_%s",i,s)
}