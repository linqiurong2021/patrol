package services

import (
	"context"
	"github.com/linqiurong2021/patrol/src/services/patrol/middleware"
	"github.com/linqiurong2021/patrol/src/services/patrol/server"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
)

type Service struct {
	RpcCli *server.GrpcClient
}

//
func NewService() *Service  {
	// redis

	return &Service{
		RpcCli: server.NewGrpcClient(),
	}
}
// 获取
func (s *Service) GetPatrol(ctx context.Context ,ID int64) (*structs.Patrol, error)  {
	//
	patrol, err := s.RpcCli.GetPatrol(ctx,ID)
	if err != nil {
		return nil, err
	}

	return patrol, nil
}

// 登录
func (s *Service) GetPatrolList(ctx context.Context ,list *structs.GetPatrolListRequest) (*structs.GetPatrolListResponse, error) {
	//
	return s.RpcCli.GetPatrolList(ctx,list)
}


// 登录
func (s *Service) Post(ctx context.Context ,postData *structs.PostRequest) (*structs.Patrol, error) {
	//
	return s.RpcCli.PostPatrol(ctx,&structs.Patrol{Memo: postData.Memo,RequestID: postData.RequestID})
}



func (s *Service) GetUserID(Token string) int64  {
	//
	middleware.Redis.HGet(Token,"id")
	return 1
}






