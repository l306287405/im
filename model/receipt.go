package model

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"im/service/cache"
	"im/service/orm"
	"strconv"
	"sync"
	"time"
)

type Receipts struct {
	Mid  uint64 `json:"mid"`
	Type byte   `json:"type"`	//0是单聊 1是群聊
}

func (r Receipts) Marshal() []byte{
	b,_:=json.Marshal(r)
	return b
}
func (r Receipts) UnMarshal(v []byte) {
	json.Unmarshal(v,r)
}

type chatReceipts struct {
	cache  *redis.Client	`json:"-"`
	db     *xorm.Engine		`json:"-"`
}

const(
	SENT = iota
	READ = iota
)

var (
	chatInstance *chatReceipts
	once sync.Once
)

func ChatReceipts() *chatReceipts {

	once.Do(func() {
		chatInstance = &chatReceipts{cache:cache.Init(),db:orm.GetDB()}
	})
	return chatInstance
}

func (r *chatReceipts) keyGenerate(mid uint64,from uint64,to uint64) string{
	return strconv.Itoa(cache.MSG_RECEIPT_MAP)+"_"+strconv.FormatUint(mid,10)+"_"+strconv.FormatUint(from,10)+"_"+strconv.FormatUint(to,10)
}

func (r *chatReceipts) Add(mid uint64,from uint64,to uint64) bool{
	err:=r.cache.Set(r.keyGenerate(mid,from,to),SENT,time.Hour).Err()
	if err!=nil{
		return false
	}
	return true
}

func (r *chatReceipts) Update(mid uint64,from uint64,to uint64) bool{

	return true
}

func (r *chatReceipts) Status(mid uint64,from uint64,to uint64) (int,error){
	return r.cache.Get(r.keyGenerate(mid,from,to)).Int()
}