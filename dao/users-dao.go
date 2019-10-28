package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"im/model"
	"im/service/cache"
	"im/service/orm"
	"strconv"
)

type UsersDao struct {
	db *xorm.Engine
}

func NewUsersDao() *UsersDao{
	return &UsersDao{db:orm.GetDB()}
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

//用户信息获取
func (u *UsersDao) Info(userId uint64) *model.Users{
	user:=new(model.Users)
	find,err:=u.db.Where("id=?",userId).Get(user)
	if err!=nil{
		panic(err)
	}
	if !find{
		return nil
	}
	return user
}

func (u *UsersDao) UpdateById(userId uint64,data *model.Users,fields ...string) error{
	_,err:=u.db.Cols(fields...).Where("id=?",userId).Update(data)
	return err
}