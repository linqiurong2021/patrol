package main

import (
	"github.com/linqiurong2021/patrol/src/services/patrol/server"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"log"
)
func initConf()  {
	confFile := "../../conf/conf.yaml"
	err := conf.NewConf(confFile)
	if err !=nil {
		log.Fatalf("init conf unmarshal error: %s",err)
	}
}
func main()  {
	initConf()
	server.NewGrpcServer()
}
