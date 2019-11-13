package dao

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"im/service/cache"
	"im/service/orm"
	"strconv"
	"sync"
	"time"
)

type chatReceiptsDao struct {
	cache  *redis.Client	`json:"-"`
	db     *xorm.Engine		`json:"-"`
}

const(
	SENT = iota
	READ = iota
)

var (
	chatInstance *chatReceiptsDao
	once sync.Once
)

func ChatReceiptsDao() *chatReceiptsDao {

	once.Do(func() {
		chatInstance = &chatReceiptsDao{cache:cache.Init(),db:orm.GetDB()}
	})
	return chatInstance
}

func (r *chatReceiptsDao) keyGenerate(mid uint64,from uint64,to uint64) string{
	return strconv.Itoa(cache.MSG_RECEIPT_MAP)+"_"+strconv.FormatUint(mid,10)+"_"+strconv.FormatUint(from,10)+"_"+strconv.FormatUint(to,10)
}

func (r *chatReceiptsDao) Add(mid uint64,from uint64,to uint64) bool{
	err:=r.cache.Set(r.keyGenerate(mid,from,to),SENT,time.Hour).Err()
	if err!=nil{
		return false
	}
	return true
}

func (r *chatReceiptsDao) Update(mid uint64,from uint64,to uint64) bool{

	return true
}

func (r *chatReceiptsDao) Status(mid uint64,from uint64,to uint64) (int,error){
	return r.cache.Get(r.keyGenerate(mid,from,to)).Int()
}
