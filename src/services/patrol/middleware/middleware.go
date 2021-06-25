package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linqiurong2021/patrol/src/services/patrol/libs"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"log"
	"strconv"
)

var Redis *libs.Redigo




// token 校验
func TokenValid(Token string) (int64,error)  {
	//
	redisAddress := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	Redis = libs.NewRedigo(redisAddress,conf.Config.Redis.Password,conf.Config.Redis.Selected )
	log.Printf("NewGrpcServer#redisAddress#%s\n",redisAddress)
	idStr := Redis.HGet(Token,"id")
	log.Printf("redis get id %s\n", idStr)
	return strconv.ParseInt(idStr,10,64)
}

//
func Auth() gin.HandlerFunc {



	// 后期使用 jwtToken
	return func(c *gin.Context) {
		//
		// 获取token
		token := c.Request.Header.Get("token")
		log.Printf("get Token %s\n",token)
		// 判断token是否合法 如果不合法则直接返回
		// 如果不合法则abort 否则 next
		id, err := TokenValid(token)
		if  err !=nil {
			result := libs.Unauthorized("token invalidate", nil)
			c.JSON(401, result)
			c.Abort()
		}
		// 可以设置request_id 调用链接监控
		// 设置user_id
		c.Set("user_id", id)
		// before request
		c.Next()
	}
}