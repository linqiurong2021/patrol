package libs

import (
	"github.com/go-redis/redis"
	"time"
)

type Redigo struct {
	Conn *redis.Client
}

//
func NewRedigo(address string,password string, selected int) *Redigo  {
	//address := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	conn := redis.NewClient(&redis.Options{
		Addr:address,
		Password: password,
		DB: selected,
	})
	return &Redigo{
		Conn: conn,
	}
}

// 存储
func (r *Redigo) GeoAdd (key string,location *redis.GeoLocation)  {
	
	r.Conn.GeoAdd(key, location)
}
func (r *Redigo) GetPoint (key string,member string)  []*redis.GeoPos {
	
	return r.Conn.GeoPos(key, member).Val()
}
//
func (r *Redigo) Ping()  {
	
	r.Conn.Ping()
}

func (r *Redigo) GetDist (key string,member1 string,member2 string)  {
	
	r.Conn.GeoDist(key, member1,member2,"m");
}
func (r *Redigo) GeoRadius (key string,longitude float64,latitude float64,radius float64)  {
	
	query := &redis.GeoRadiusQuery{
		Radius: radius,
		Unit: "m",
	}
	r.Conn.GeoRadius(key, longitude,latitude,query);
}


//
func (r *Redigo) Set(key string,value interface{},expireAt time.Duration)  {
	
	r.Conn.Set(key,value,expireAt)
}

//
func (r *Redigo) Get(key string)  {
	
	r.Conn.Get(key)
}

func (r *Redigo) Del(key string)  {
	
	r.Conn.Del(key)
}

func (r *Redigo) SetExpireAt(key string,expireAt time.Duration)  {
	
	r.Conn.Expire(key,expireAt)
}

func (r *Redigo) ZRange(key string,start,end int64) []string {
	
	return r.Conn.ZRange(key,start,end).Val()
}

func (r *Redigo) ZCount(key string,min string,max string) int64  {
	
	return r.Conn.ZCount(key,min,max).Val()
}
//
func (r *Redigo) ZCountAll(key string) int64 {

	return r.Conn.ZCount(key,"-inf","+inf").Val()
}

func (r *Redigo) ZMembers(key string) []string {

	zCount := r.Conn.ZCount(key, "-inf","+inf").Val()
	zRange := r.Conn.ZRange(key,0,zCount).Val()
	return zRange
}


//
func (r *Redigo) ZPoints(key string, members []string) [][]*redis.GeoPos  {
	//
	var points [][]*redis.GeoPos
	for _, member := range members {
		//
		point := r.Conn.GeoPos(key, member).Val()
		points = append(points, point)
	}
	return points
}

//
func (r *Redigo) HSet(key string, hKey string,hValue interface{})  {
	
    r.Conn.HSet(key,hKey,hValue)
}

//
func (r *Redigo) HGet(key string, hKey string) string  {
	
	strCmd := r.Conn.HGet(key,hKey)
	return strCmd.Val()
}



