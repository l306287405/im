package service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/iris-contrib/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
	"im/model"
	"im/service/cache"
	"im/service/orm"
	"os"
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
func (s *UserService) Login(appId uint,username string,requestPassword []byte) (*string,error){
	var(
		token *string
		err error
	)

	has,err := s.db.Cols("id","apps_id","account","nickname","token","password",).Where("apps_id=?",appId).
		And("account=?",username).And("status=?",1).Get(&s.users)
	if err!=nil{
		return token,err
	}
	if !has{
		return token,errors.New("账号或密码错误")
	}
	passwordRecord := []byte(s.users.Password)
	err=bcrypt.CompareHashAndPassword(passwordRecord,requestPassword)
	if err!=nil{
		return token,errors.New("账号或密码错误")
	}

	token,err=s.CreateToken()
	if err!=nil{
		return token,err
	}

	err=s.SetCacheOfToken(*token,s.users.Id)
	if err!=nil{
		return token,errors.New("用户token存储失败 请联系管理员 "+err.Error())
	}

	return token,nil
}

//创建token
func (s *UserService) CreateToken() (*string,error){
	token:=jwt.NewTokenWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"id":s.users.Id,
		"apps_id":s.users.AppsId,
		"account":s.users.Account,
		"nickname":s.users.Nickname,
		"exp":time.Now().Add(7*24*time.Hour).Unix(),
	})
	tokenString,err:=token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	s.users.Token=&tokenString
	return s.users.Token,err
}

//创建一个账号
func (s *UserService) Create(appsId uint,account string,password string,nickname string) (result int64,err error){
	p, _ :=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	password = string(p)
	s.users.AppsId,s.users.Account,s.users.Password,s.users.Nickname,s.users.Status=appsId,account,password,nickname,1

	return s.db.InsertOne(&s.users)
}

//设置token与用户id的映射缓存
func (s *UserService) SetCacheOfToken(token string,userId uint64) error{
	return cache.Init().HSet(string(cache.USERS_TOKEN_MAP),string(userId),token).Err()
}

//获取token相关的用户id
func (s *UserService) GetCacheByUid(userId uint64) (string,error){
	jwt,err:=cache.Init().HGet(string(cache.USERS_TOKEN_MAP),string(userId)).Result()
	if err!=nil{
		return "0",err
	}
	return jwt,nil
}

//删除token相关的缓存
func (s *UserService) DelCacheOfToken(userId uint64) error{
	return cache.Init().HDel(string(cache.USERS_TOKEN_MAP),string(userId)).Err()
}