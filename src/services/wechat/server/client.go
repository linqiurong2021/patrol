package server

import (
	"context"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
	"google.golang.org/grpc"
	"log"
	pb "github.com/linqiurong2021/patrol/src/services/wechat/protof/v1"
)



type GrpcClient struct {
	Client pb.WechatServiceClient
	Connect *grpc.ClientConn
}

func NewGrpcClient() *GrpcClient{
	// 端口号
	address := conf.Config.Grpc.Port
	// 拨号
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewWechatServiceClient(conn)
	return  &GrpcClient{
		Client: client,
		Connect: conn,
	}
}

//
func (c *GrpcClient) Login(ctx context.Context ,user *structs.Login)  (*structs.Login, error) {

	reqUser :=  pb.User{
		Nickname: user.Nickname,
		OpenID: user.OpenID,
		Gender: user.Gender,
		City: user.City,
		Province: user.Province,
		Country: user.Country,
		AvatarUrl: user.AvatarUrl,
	}
	fmt.Printf("#login#\n %#v #reqUser \n",&reqUser)
	// 请求数据
	req := pb.LoginRequest{
		User: &reqUser,
	}
	log.Printf("Login req params %#v\n",&req)
	//
	resp ,err := c.Client.Login(ctx,&req)
	//
	fmt.Sprintf("%#v,#########",resp)
	if err != nil {
		log.Printf("client login error %s\n", err)
		return nil ,err
	}
	return &structs.Login{
		Nickname: resp.User.Nickname,
		OpenID: resp.User.OpenID,
		Gender: resp.User.Gender,
		City: resp.User.City,
		Province: resp.User.Province,
		Country: resp.User.Country,
		AvatarUrl: resp.User.AvatarUrl,
		Token: resp.User.Token,
	},nil
}

func (c *GrpcClient) Code2Session(ctx context.Context , code string) (*structs.Code2Session, error)  {
	// 请求数据
	req := pb.Code2SessionRequest{
		Code: code,
	}
	//
	resp ,err := c.Client.Code2Session(ctx,&req)
	//
	fmt.Sprintf("%#v,#########",resp)
	if err != nil {
		fmt.Printf("client code2session error %s\n", err)
		//
	}
	// 如果不为0 则说明是微信调用出错
	if resp.ErrCode !=0 {
		fmt.Printf("微信调用错误了 error %s", resp.ErrMsg)
	}
	//
	fmt.Printf("OPENID: %s,SESSION_KEY:%s\n",resp.OpenID,resp.SessionKey)

	return &structs.Code2Session{SessionKey: resp.SessionKey,OpenID: resp.OpenID}, err
}

func GetUser()  {
	
}

