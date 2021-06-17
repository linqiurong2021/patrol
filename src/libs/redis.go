package libs

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"log"
)

type Redis struct {
	Conn redis.Conn
}

//
func NewRedis() *Redis  {
	address := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	conn, err := redis.Dial("tcp",address)
	if err != nil {
		log.Fatalf("connect redis error: %s", err)
	}
	return &Redis{
		Conn: conn,
	}
}

// 存储
func (r *Redis) Store ()  {
	defer r.Conn.Close()
	_, err := r.Conn.Do("geoadd")
	if err != nil {
		//
	}
}
