package server

import (
	"encoding/json"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/consts"
	"github.com/linqiurong2021/patrol/src/services/wechat/libs"
	"github.com/linqiurong2021/patrol/src/services/wechat/model"
	pb "github.com/linqiurong2021/patrol/src/services/wechat/protof/v1"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type GrpcServer struct {
	Model *model.Model
	Redis *libs.Redis
}


func (g *GrpcServer) existsUser(OpenID string) (*model.User, error)  {
	user, err := g.Model.GetUserByOpenID(OpenID)
	if err != nil {
		return nil, err
	}
	return user, nil;
}

func (g *GrpcServer) generateToken(data *model.User) string {
	// 生成token
	return libs.NewToken(data.Token);
}


func (g *GrpcServer) cache(token string, data *model.User)  {
	//
	g.Redis.HSet(token,"id",data.ID)
}

func (g *GrpcServer) GetUserID(token string)   {
	g.Redis.HGet(token,"id")
}
// 登录
func (g *GrpcServer) doLogin(req *pb.LoginRequest) (*model.User, error) {
	log.Printf("#server#doLogin\n")
	// 判断用户是否存在 如果不存在则更新 如果存在则更新
	user, err := g.existsUser(req.User.OpenID)
	// 出错时
	if err!=nil {
		return nil,err
	}
	log.Printf("#user#: %#v\n,req.User %#v",user,req.User)
	// 无记录时
	newUser := model.User{
		Nickname: req.User.Nickname,
		OpenID: req.User.OpenID,
		Gender: req.User.Gender,
		City: req.User.City,
		Province: req.User.Province,
		Country: req.User.Country,
		AvatarUrl: req.User.AvatarUrl,
	}

	if user.ID == 0 {
		fmt.Print("Create")
		log.Printf("#create before data %#v\n", &newUser)
		// 创建用户并生成token
		user, err = g.Model.CreateUser(&newUser)
		if err != nil {
			return nil ,err
		}
		// 创建token
		user.Token = g.generateToken(user)

		return user, nil
	}else{
		// 有记录时,更新数据
		newUser.ID = user.ID
		_, err := g.Model.UpdateUser(&newUser)
		if err != nil {
			log.Printf("update user error %#v\n", err)
			return nil, err
		}

		log.Print("update user success \n")
		//
		user, err = g.Model.GetUser(user.ID)
		if err != nil {
			return nil ,err
		}
		user.Token = g.generateToken(user)
		log.Printf("return user %#v\n",user)
		return user,nil
	}
}

//
func (g *GrpcServer)Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 登录操作
	log.Printf("#server#Login\n")
	user, err := g.doLogin(in)
	if err != nil {
		return nil, err
	}
	// 登录返回
	return  &pb.LoginResponse{User: &pb.User{
		Nickname: user.Nickname,
		OpenID: user.OpenID,
		Gender: user.Gender,
		City: user.City,
		Province: user.Province,
		Country: user.Country,
		AvatarUrl: user.AvatarUrl,
		Token: user.Token,
	}},nil
}

func (g *GrpcServer) GetWxResponse(appID string,secret string,Code string) (*http.Response, error)  {
	//
	log.Printf("APPID: %s, Secret: %s, Code: %s\n",appID,secret,Code)
	//
	URL := consts.GetCode2SessionURL(appID,secret,Code)
	//
	log.Printf("Code2Session URL %s\n",URL)
	//
	return http.Get(URL)
}
//
func (g *GrpcServer) Code2Session(ctx context.Context,req *pb.Code2SessionRequest) (*pb.Code2SessionResponse, error)  {
	//
	code := req.Code
	log.Printf("%s,request.code",code)

	appID := "wxf5f6d29c7f7b73a9"
	secret := "6901d6d93a764abed06118dd608a9ade"
	//
	response , err := g.GetWxResponse(appID,secret,code)
	log.Printf("%#v,response.Body",response.Body)
	if err != nil {
		fmt.Printf("get code2session body error %s", err)
	}
	//
	defer response.Body.Close()
	//
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read body error %s", err)
	}
	log.Printf("read body :%s",string(body[:]))
	//
	var code2Session structs.Code2Session
	err = json.Unmarshal(body,&code2Session)
	log.Printf("unmashal code2session:  %#v\n", code2Session)
	//
	return &pb.Code2SessionResponse{
		ErrCode: code2Session.ErrCode,
		ErrMsg: code2Session.ErrMsg,
		OpenID: code2Session.OpenID,
		SessionKey: code2Session.SessionKey,
	}, nil
}


//
func (g *GrpcServer)GetUser(ctx context.Context, req *pb.GetUserRequest ) (*pb.GetUserResponse ,  error){
	return nil,nil
}



// GrpcServer
func NewGrpcServer()  {
	//
	address := ":8081" //conf.Config.Grpc.Port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//
	redisAddress := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	redis := libs.NewRedis(redisAddress,conf.Config.Redis.Password,conf.Config.Redis.Selected )
	log.Printf("NewGrpcServer#redisAddress#%s\n",redisAddress)
	// postgres dataSourceName
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName)
	log.Printf("NewGrpcServer#dataSourceName#%s\n",dataSourceName)
	pb.RegisterWechatServiceServer(s, &GrpcServer{
		Model: model.NewModel(dataSourceName),
		Redis: redis,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Printf("#Wechat#GrpcServer#Start At# %s",address)

}


