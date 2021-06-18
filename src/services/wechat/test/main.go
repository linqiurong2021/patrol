package main

import (
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/server"
	"log"
)
func init()  {
	// 配置文件初始化
	confFile := "../conf/conf.yaml"
	err := conf.NewConf(confFile)
	if err !=nil {
		log.Fatalf("init conf unmarshal error: %s",err)
	}
}
func main()  {
	server.NewGrpcServer()
}


