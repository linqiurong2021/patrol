package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/libs"
	"github.com/linqiurong2021/patrol/src/services/wechat/services"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
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
	//go func() {
	//	server.NewGrpcServer2()
	//}()
	gin.Run(address)
}

// code 2 session
func (h *Handler) Code2Session(c *gin.Context)  {
	code := c.Query("code")
	if code == "" {
		// 判断校验
		result := libs.ValidFailure("code must");
		c.JSON(200,result)
	}
	log.Printf("CODE:%s\n",code)
	resp, err := h.Srv.Code2Session(c,code)
	if err !=nil{
		msg := fmt.Sprintf("srv.Code2Session error %s",err)
		result := libs.ServerError(msg,nil)
		c.JSON(400,result)
	}
	result := libs.Success("success",resp)
	c.JSON(200,result)
}
// login
func (h *Handler) Login(c *gin.Context)  {
	//
	var login structs.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		result := libs.ValidFailure(nil)
		c.JSON(200,result)
		return
	}
	fmt.Printf("%#v,request body \n",login)
	resp, err :=h.Srv.Login(c,&login)

	if err !=nil{
		msg := fmt.Sprintf("srv.Login error %s",err)
		result := libs.ServerError(msg,nil)
		c.JSON(400,result)
		return
	}
	result := libs.Success("success",resp)
	c.JSON(200,result)
}

// getUser
func (h *Handler) GetUser(ctx *gin.Context)  {
	
}
