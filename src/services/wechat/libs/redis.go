package libs

import (
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	Conn *redis.Client
}

//
func NewRedis(address string,password string, selected int) *Redis  {
	//address := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	conn := redis.NewClient(&redis.Options{
		Addr:address,
		Password: password,
		DB: selected,
	})
	return &Redis{
		Conn: conn,
	}
}

// 存储
func (r *Redis) GeoAdd (key string,location *redis.GeoLocation)  {
	defer r.Conn.Close()
	r.Conn.GeoAdd(key, location)
}
func (r *Redis) GetPoint (key string,member string)  {
	defer r.Conn.Close()
	r.Conn.GeoPos(key, member)
}
//
func (r *Redis) Ping()  {
	r.Conn.Ping()
}

func (r *Redis) GetDist (key string,member1 string,member2 string)  {
	defer r.Conn.Close()
	r.Conn.GeoDist(key, member1,member2,"m");
}
func (r *Redis) GeoRadius (key string,longitude float64,latitude float64,radius float64)  {
	defer r.Conn.Close()
	query := &redis.GeoRadiusQuery{
		Radius: radius,
		Unit: "m",
	}
	r.Conn.GeoRadius(key, longitude,latitude,query);
}


//
func (r *Redis) Set(key string,value interface{},expireAt time.Duration)  {
	defer r.Conn.Close()
	r.Conn.Set(key,value,expireAt)
}

//
func (r *Redis) Get(key string)  {
	defer r.Conn.Close()
	r.Conn.Get(key)
}

func (r *Redis) Del(key string)  {
	r.Conn.Del(key)
}

func (r *Redis) SetExpireAt(key string,expireAt time.Duration)  {
	r.Conn.Expire(key,expireAt)
}

func (r *Redis) ZRange(key string,start,end int64)  {
	r.Conn.ZRange(key,start,end)
}

func (r *Redis) ZCount(key string,min string,max string)  {
	r.Conn.ZCount(key,min,max)
}
//
func (r *Redis) ZCountAll(key string)  {
	r.ZCount(key,"-inf","+inf")
}



//
func (r *Redis) HSet(key string, hKey string,hValue interface{}) bool  {
	defer r.Conn.Close()
    boolCmd := r.Conn.HSet(key,hKey,hValue)
    return boolCmd.Val()
}

//
func (r *Redis) HGet(key string, hKey string)  {
	defer r.Conn.Close()
	r.Conn.HGet(key,hKey)
}



