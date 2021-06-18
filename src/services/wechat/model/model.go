package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"log"
)

type Model struct {
	engine *xorm.Engine
}

func init()  {
	NewModel()
}
//
func NewModel() *Model {
	// 连接字符串
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName) // db
	engine, err := xorm.NewEngine("postgres", dataSourceName)
	if err != nil {
		//
		log.Fatalf("init db error: %s", err)
	}
	return &Model{
		engine: engine,
	}
}

// 用户表
type User struct {
	ID       uint64 `json:"id";xorm:"id"`
	OpenID 	string  `json:"openid";xorm:"openid"`
	NickName string `json:"nickname";xorm:"nickname"`
	Avatar   string `json:"avatar";xorm:"avatar"`
	Gender   int8   `json:"gender";xorm:"gender"`
	Country  string `json:"country";xorm:"country"`
	Province string `json:"province";xorm:"province"`
	City     string `json:"city";xorm:"city"`
	Token 	string `json:"token";`
}

// 获取用户信息
func (m *Model) GetUser(id uint64) (user *User, err error) {
	_, err = m.engine.Table(&user).Where("id = ?", id).Cols("id,openid,name,province,gender,city,country,avatar").Get(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 获取用户信息
func (m *Model) GetUserByOpenID(OpenID string) (user *User, err error) {
	_, err = m.engine.Table(&User{}).Where("openid = ?", OpenID).Cols("id,openid,name,province,gender,city,country,avatar").Get(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}


// 创建用户
func (m *Model) CreateUser(user *User) (outUser *User, err error) {
	affected, err := m.engine.Table(&User{}).Insert(&user)
	if err != nil {
		return nil, err
	}
	if affected > 0 {
		return user, nil
	}
	return nil, nil
}

// 更新用户
func (m *Model) UpdateUser(user *User) (affected int64, err error) {
	return m.engine.Table(&User{}).Where("openid = ?", user.OpenID).Update(&user)
}
