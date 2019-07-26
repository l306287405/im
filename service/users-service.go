package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
	"im/model"
	"im/service/orm"
	"math/rand"
	"time"
)

type UserService struct {
	db *xorm.Engine
	users model.Users
}

func NewUserService() *UserService{
	return &UserService{db:orm.GetDB()}
}

//登录模块
func (us *UserService) Login(username string,requestPassword []byte) (*string,error){
	var(
		token *string
		err error
	)

	has,err := us.db.Cols("id","token","password",).
		Where("account=?",username).And("status=?",1).Get(&us.users)
	if err!=nil{
		return token,err
	}
	if !has{
		return token,errors.New("账号或密码错误")
	}
	passwordRecord := []byte(us.users.Password)
	err=bcrypt.CompareHashAndPassword(passwordRecord,requestPassword)
	if err!=nil{
		return token,errors.New("账号或密码错误")
	}
	token,err=us.CreateToken(us.users.Id)
	if err!=nil{
		return token,err
	}
	return token,nil
}

//创建token
func (us *UserService) CreateToken(userId uint64) (*string,error){
	build := []byte(fmt.Sprintf("%d%s%d",userId,time.Now(),rand.Int()))
	hash := fmt.Sprintf("%x",sha256.Sum256(build))
	us.users.Token=&hash
	_,err := us.db.Cols("token").ID(userId).Update(&us.users)
	return us.users.Token,err
}

//创建一个账号
func (us *UserService) Create(account string,password string,nickname string) (result int64,err error){
	p, _ :=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	password = string(p)
	us.users.Account,us.users.Password,us.users.Nickname,us.users.Status=account,password,nickname,1


	return us.db.InsertOne(&us.users)
}