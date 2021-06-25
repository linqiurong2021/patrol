package test

import (
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"github.com/linqiurong2021/patrol/src/services/patrol/libs"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
	"testing"
)

func TestMain(m *testing.M)  {
	initConf()
	m.Run()
}


func TestZCountAll(t *testing.T)  {
	redisAddress := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,6379)
	fmt.Printf("redisAddress:%s\n", redisAddress)

	redis := libs.NewRedigo(redisAddress,conf.Config.Redis.Password,conf.Config.Redis.Selected )
	// docker读取不到
	key := "7f550786c810e0ea5ec97446439156f1"
	//
	defer redis.Conn.Close()
	//
	members := redis.ZMembers(key)
	//
	zPoints := redis.ZPoints(key, members)
	//
	var points []*structs.Point
	for _,geoPos := range zPoints {
		for _,item := range geoPos {
			points = append(points, &structs.Point{
				Longitude: item.Longitude,
				Latitude: item.Latitude,
			})
		}
	}
	//
	for _,point := range points{
		fmt.Printf("point: %#v\n",point)
	}


	fmt.Printf("members %#v\n",points)
	//defer redis.Conn.Close()
	//key := "7f550786c810e0ea5ec97446439156f1"
	//zcount := redis.Conn.ZCount(key,"-inf","+inf").Val()
	//val := redis.Conn.ZRange(key,0, zcount).Val()
	////////
	//fmt.Printf("zrange %#v\n val:%d\n",val,zcount)
	////fmt.Printf("count: %#v, zrange: %#v\n, max: %#v, list: %#v \n",count,zrange, max, list)
}
