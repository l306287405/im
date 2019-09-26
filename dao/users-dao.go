package dao

import (
	"fmt"
	"im/service/cache"
	"strconv"
)

type UsersDao struct {

}

func NewUsersDao() *UsersDao{
	return &UsersDao{}
}

//获取用户上线信息的cache_key
func onlineKey(appId uint,uid uint64) string{
	return fmt.Sprintf("%d_%d",appId,uid)
}

//用户上线
func (u *UsersDao) Online(appId uint,uid uint64,cid string) error{
	return cache.Init().HSet(strconv.Itoa(cache.USERS_COMM_MAP),onlineKey(appId,uid),cid).Err()
}

//用户下线
func (u *UsersDao) OffLine(appId uint,uid uint64) error{
	return cache.Init().HDel(strconv.Itoa(cache.USERS_COMM_MAP),onlineKey(appId,uid)).Err()
}

//