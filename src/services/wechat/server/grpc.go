package server

import (
	"encoding/json"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/wechat/consts"
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
}


//
func (g *GrpcServer)Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	return  nil,nil
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
	pb.RegisterWechatServiceServer(s, &GrpcServer{});

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Printf("#Wechat#GrpcServer#Start At# %s",address)

}


