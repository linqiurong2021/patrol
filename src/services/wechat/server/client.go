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
	fmt.Printf("dial grpc server error %s",err)
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
func Login()  {
	//
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

