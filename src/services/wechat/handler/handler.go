package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/services"
	"log"
)

//
type Handler struct {
	Srv *services.Service
}


func NewHandler() *Handler {
	return &Handler{
		Srv: services.NewService(),
	}
}

func RunHttpServer() {
	//
	gin := gin.Default()
	// 版本1
	v1 := gin.Group("/v1/")
	{
		//
		handler := NewHandler()
		//
		v1.GET("/code2session",handler.Code2Session)
		//
		v1.POST("/login", handler.Login)
		//
		v1.GET("/getUser", handler.GetUser)
	}
	address := conf.Config.App.Port
	fmt.Printf("address %s", address)
	//
	gin.Run(address)
}

// code 2 session
func (h *Handler) Code2Session(c *gin.Context)  {
	code := c.Query("code")
	if code == "" {
		// 判断校验
		c.JSON(400,"code must")
	}
	log.Printf("CODE:%s\n",code)
	resp, err := h.Srv.Code2Session(c,code)
	if err !=nil{
		fmt.Printf("srv.Code2Session error %s", err)
	}
	c.JSON(200, resp)
}
// login
func (h *Handler) Login(ctx *gin.Context)  {
	//
}

// getUser
func (h *Handler) GetUser(ctx *gin.Context)  {
	
}
