package member

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"im/model"
	"im/service/orm"
	"math/rand"
	"time"
)

//登录模块
func Login(username string,requestPassword []byte) (*string,error){
	var(
		db=orm.GetDB()
		users=new(model.Users)
		token *string
		err error
	)

	has,err := db.Cols("id","token","password",).
		Where("account=?",username).And("status=?",1).Get(users)
	if !has{
		return token,errors.New("账号或密码错误")
	}
	passwordRecord := []byte(users.Password)
	err=bcrypt.CompareHashAndPassword(passwordRecord,[]byte(requestPassword))
	if err!=nil{
		return token,errors.New("账号或密码错误")
	}
	token,err=CreateToken(users.Id)
	if err!=nil{
		return token,err
	}
	return token,nil
}

//创建token
func CreateToken(userId uint64) (*string,error){
	build := []byte(fmt.Sprintf("%d%s%d",userId,time.Now(),rand.Int()))
	hash := fmt.Sprintf("%x",sha256.Sum256(build))
	users := new(model.Users)

	users.Token=&hash
	_,err := orm.GetDB().Cols("token").ID(userId).Update(users)
	if err!=nil{
		return users.Token,err
	}
	return users.Token,nil
}

//创建一个账号
func Create(account string,password string,nickname string) (result int64,err error){
	var(
		db = orm.GetDB()
		users = model.Users{Account:account,Nickname:nickname,Status:1}
	)

	p, _ :=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	password = string(p)
	users.Password=password

	return db.InsertOne(users)
}