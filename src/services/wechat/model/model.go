package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

type Model struct {
	engine *xorm.Engine
}

//
func NewModel(dataSourceName string) *Model {
	// 连接字符串
	//dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName) // db
	log.Printf("### dataSourceName:%s\n",dataSourceName)
	engine, err := xorm.NewEngine("postgres", dataSourceName)
	if err != nil {
		//
		log.Fatalf("init db error: %s", err)
	}
	err = engine.Sync2(new(User))
	if err != nil {
		log.Fatalf("engine sync error %s\n",err)
	}
	return &Model{
		engine: engine,
	}
}

// 用户表
type User struct {
	ID       int64 `json:"id" xorm:"id pk autoincr"` // 不要写 xorm:"id" 否则插入会提示  cannot insert into column "id"
	OpenID 	string  `json:"openid" xorm:"openid"`
	Nickname string `json:"nickname" xorm:"nickname"`
	AvatarUrl   string `json:"avatar_url" xorm:"avatar_url"`
	Gender   uint32   `json:"gender" xorm:"gender"`
	Country  string `json:"country" xorm:"country"`
	Province string `json:"province" xorm:"province"`
	City     string `json:"city" xorm:"city"`
	Token 	string `json:"token" `
}

// 获取用户信息
func (m *Model) GetUser(id int64) ( *User,  error) {
	var user User;
	_, err := m.engine.Table(&User{}).Where("id = ?", id).Cols("id,openid,nickname,province,gender,city,country,avatar_url").Get(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 获取用户信息
func (m *Model) GetUserByOpenID(OpenID string) (*User,  error) {
	var user User
	has, err := m.engine.Table(&User{}).Where("openid = ?", OpenID).Cols("id,openid,nickname,province,gender,city,country,avatar_url").Get(&user)
	fmt.Printf("get record %#v\n, %b\n",err,has)
	if err != nil && err != xorm.ErrNotExist{
		return nil, err
	}
	return &user, nil
}


// 创建用户
func (m *Model) CreateUser(userData *User) (outUser *User, err error) {
	log.Printf("insert data %#v\n",userData)
	affected, err := m.engine.Insert(userData)

	if err != nil {
		log.Printf("insert data error %s\n", err)
		return nil, err
	}
	if affected > 0 {
		userData.ID = affected
		return userData, nil
	}
	return nil, nil
}

// 更新用户
func (m *Model) UpdateUser(user *User) (affected int64, err error) {
	log.Printf("update user data %#v\n",user)
	return m.engine.Where("id = ?",user.ID).AllCols().Omit("id").Update(user)
}
