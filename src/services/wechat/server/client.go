package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	pb "github.com/linqiurong2021/patrol/src/services/wechat/protobuf/v1"
	"log"
	"time"
)

type Etcd struct {
	Client *clientv3.Client
}
//
func NewEtcd(endpoints []string)  *Etcd{
	//

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatalf("init etcd client error %s\n", err)
	}
	return &Etcd{
		Client: cli,
	}
}

func NewGrpcClient2() *GrpcClient{
	cli, err := clientv3.NewFromURL(conf.Config.Register.Host)
	if err !=nil {
		log.Fatalf("error client v3 newFormURL %s\n", err)
	}
	// cli
	serviceKey := fmt.Sprintf("%s/%s",conf.Config.Grpc.Name,conf.Config.Grpc.Port)
	fmt.Printf("service Key: %s\n",serviceKey)
	resp , err := cli.Get(context.TODO(),serviceKey)
	//
	var addr []string
	var service Service
	for _,item := range resp.Kvs{
		// 负载均衡
		err =json.Unmarshal([]byte(item.Value),&service)
		if err != nil {
			log.Printf("json unmashal err %s\n", err)
		}
		addr = append(addr, service.Addr)
	}
	//
	fmt.Printf("serverAddr: %#v\n",addr[0])
	// 拨号
	conn, err := grpc.Dial(addr[0], grpc.WithInsecure() )
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	client := pb.NewWechatServiceClient(conn)
	return  &GrpcClient{
		Client: client,
		Connect: conn,
		//Etcd: etcd,
	}
	//return nil
}

type GrpcClient struct {
	Client pb.WechatServiceClient
	Connect *grpc.ClientConn
	Etcd *Etcd
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
		//Etcd: etcd,
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

