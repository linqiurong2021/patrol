package test

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"log"
	"testing"
	"time"
)

type CreateTime struct {
	ID int64
	CreateAt time.Time `xorm:"create_at"`
}

func TestCreateAt(t *testing.T)  {
	//
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName)
	log.Printf("NewGrpcServer#dataSourceName#%s\n",dataSourceName)
	engine, err := xorm.NewEngine("postgres", dataSourceName)
	if err != nil {
		fmt.Printf("#error %s \n ", err)
	}
	engine.Sync2(&CreateTime{})
	now := time.Now()
	_ ,err =engine.Exec("insert into create_time(create_at) values(?)",now)
	if err !=nil {
		fmt.Printf("exec error %s\n", err)
	}
}


