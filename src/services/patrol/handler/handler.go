package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"github.com/linqiurong2021/patrol/src/services/patrol/libs"
	"github.com/linqiurong2021/patrol/src/services/patrol/middleware"
	"github.com/linqiurong2021/patrol/src/services/patrol/server"
	"github.com/linqiurong2021/patrol/src/services/patrol/services"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
	"io/ioutil"
	"log"
	"strconv"
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
		// 调用前需要判断是否已登录
		v1.Use(middleware.Auth())
		//
		handler := NewHandler()
		//
		v1.GET("/getPatrol",handler.GetPatrol)
		//
		v1.GET("/getPatrolList", handler.GetPatrolList)
		//
		v1.POST("/post", handler.Post)
	}
	address := conf.Config.App.Port
	fmt.Printf("address %s", address)
	//
	go func() {
		server.NewGrpcServer()
	}()
	gin.Run(address)
}

// code 2 session
func (h *Handler) GetPatrol(c *gin.Context)  {
	//
	page := c.Query("page")
	pageSize := c.Query("page_size")
	UserID := c.MustGet("user_id")
	// 使用validator
	if page == ""  || pageSize == "" {
		// 判断校验
		result := libs.ValidFailure("page and page size must");
		c.JSON(200,result)
		return
	}
	log.Printf("page:%s pageSize:%s \n",page,pageSize)
	//
	resp, err := h.Srv.GetPatrol(c,UserID.(int64))
	if err !=nil{
		msg := fmt.Sprintf("srv.GetPatrol error %s",err)
		result := libs.ServerError(msg,nil)
		c.JSON(400,result)
		return
	}
	result := libs.Success("success",resp)
	c.JSON(200,result)
}
// login
func (h *Handler) GetPatrolList(c *gin.Context)  {
	//
	var list structs.GetPatrolListRequest
	page := c.Query("page")
	pageSize := c.Query("page_size")
	UserID := c.MustGet("user_id")
	// 使用validator
	if page == ""  || pageSize == "" {
		// 判断校验
		result := libs.ValidFailure("page and page size must");
		c.JSON(200,result)
		return
	}
	currPage, err :=strconv.ParseInt(page,10,64)
	if err != nil {
		log.Printf("page must int")
		return
	}
	currPageSize, err := strconv.ParseInt(pageSize,10,64)
	if err != nil {
		log.Printf("page_size must int")
		return
	}
	list = structs.GetPatrolListRequest{
		UserID: UserID.(int64),
		PageParam: structs.PageParam{
			Page:currPage ,
			PageSize: currPageSize,
		},
	}
	fmt.Printf("%#v,request body \n",list)
	resp, err :=h.Srv.GetPatrolList(c,&list)

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
func (h *Handler) Post(c *gin.Context)  {
	// 获取参数
	// memo request_id
	var post *structs.PostRequest
	body,_ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body,&post)

	resp, err :=h.Srv.Post(c,post)
	if err !=nil{
		msg := fmt.Sprintf("srv.Post error %s",err)
		result := libs.ServerError(msg,nil)
		c.JSON(400,result)
		return
	}
	result := libs.Success("success",resp)
	c.JSON(200,result)
}
