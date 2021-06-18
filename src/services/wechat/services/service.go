package services

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/linqiurong2021/patrol/src/services/wechat/model"
	"github.com/linqiurong2021/patrol/src/services/wechat/server"
	"github.com/linqiurong2021/patrol/src/services/wechat/structs"
)

type Service struct {
	RpcCli *server.GrpcClient
	Model *model.Model
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

// 登录
func (s *Service) Login(data *model.User) (*model.User, error) {
	// 判断用户是否存在 如果不存在则更新 如果存在则更新
	user, err := s.existsUser(data.OpenID)
	if err!=nil {
		// 查无数据时
		if err == xorm.ErrTableNotFound {
			// 创建用户并生成token
			user, err = s.Model.CreateUser(data)
			if err != nil {
				return nil ,err
			}
			// 创建token
			user.Token = s.generateToken(user)
			return user, nil
		}
		return nil,err
	}
	// 更新数据
	_, err = s.Model.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	//
	user, err = s.Model.GetUserByOpenID(data.Token)
	if err != nil {
		return nil ,err
	}
	user.Token = s.generateToken(data)
	return data,nil

}

func (s *Service) existsUser(OpenID string) (*model.User, error)  {
	user, err := s.Model.GetUserByOpenID(OpenID)
	if err != nil {
		return nil, err
	}
	return user, nil;
}

func (s *Service) generateToken(data *model.User) string {
	// 生成token
	//token := libs.NewToken(data.OpenID)
	// 存储token
	//redis := libs.NewRedis()
	//

	return "token";
}