package services

import (
	"context"
	"github.com/linqiurong2021/patrol/src/services/wechat/server"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
)

type Service struct {
	RpcCli *server.GrpcClient
}

//
func NewService() *Service  {
	// redis
	return &Service{
		RpcCli: server.NewGrpcClient2(),
	}
}
// 获取
func (s *Service) Code2Session(ctx context.Context ,code string) (*structs.Code2Session, error)  {
	return s.RpcCli.Code2Session(ctx,code)
}

// 登录
func (s *Service) Login(ctx context.Context ,login *structs.Login) (*structs.Login, error) {
	//
	return s.RpcCli.Login(ctx,login)
}
