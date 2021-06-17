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
	return &Service{
		RpcCli: server.NewGrpcClient(),
	}
}
// 获取
func (s *Service) Code2Session(ctx context.Context ,code string) (*structs.Code2Session, error)  {
	return s.RpcCli.Code2Session(ctx,code)
}