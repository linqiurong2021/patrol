package main

import (
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"github.com/linqiurong2021/patrol/src/services/patrol/handler"
	"log"
)

// 初始化
func init()  {
	// 配置文件初始化
	confFile := "./conf/conf.yaml"
	err := conf.NewConf(confFile)
	if err !=nil {
		log.Fatalf("init conf unmarshal error: %s",err)
	}
}

func main()  {
	fmt.Println("PatrolService")
	// 开启微服务协程
	handler.RunHttpServer()
}
