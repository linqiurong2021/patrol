package server

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"github.com/linqiurong2021/patrol/src/services/patrol/libs"
	"github.com/linqiurong2021/patrol/src/services/patrol/model"
	pb "github.com/linqiurong2021/patrol/src/services/patrol/protobuf/v1"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	Model *model.Model
	Redis *libs.Redigo
	Geometry *libs.Geometry
}

// 获取巡查详情
func (g *GrpcServer)GetPatrol(_ context.Context, in *pb.GetPatrolByIDRequest) (*pb.GetPatrolByIDResponse, error) {
	//
	userID := in.UserID
	ID := in.ID
	//
	patrol, err := g.Model.GetPatrol(userID,ID)
	fmt.Printf("获取patrol: %#v\n",patrol)
	if err != nil {
		return nil, err
	}
	createAt, _ := ptypes.TimestampProto(patrol.CreateAt)
	updateAt, _ := ptypes.TimestampProto(patrol.UpdateAt)
	return &pb.GetPatrolByIDResponse{
		Patrol: &pb.Patrol{
			ID: patrol.ID,
			UserID: patrol.UserID,
			Line: patrol.Polyline,
			// 转换
			CreateAt: createAt,
			UpdateAt: updateAt,
		},
	},nil
}
// 获取某用户巡查列表
func (g *GrpcServer) GetPatrolList(_ context.Context,in *pb.GetPatrolListByUserIDRequest) (*pb.GetPatrolListByUserIDResponse, error)  {
	//
	patrol, err := g.Model.GetPatrolList(in.UserID,int(in.Page),int(in.PageSize))
	if err != nil {
		return nil, err
	}
	var patrols []*pb.Patrol
	// 数据处理
	for _,item := range patrol{
		//
		createAt, _ := ptypes.TimestampProto(item.CreateAt)
		updateAt, _ := ptypes.TimestampProto(item.UpdateAt)
		patrols = append(patrols,&pb.Patrol{
			UserID: item.UserID,
			ID: item.ID,
			Line: item.Polyline,
			CreateAt:createAt ,
			UpdateAt: updateAt,
			Memo: item.Memo,
		})
	}
	count, err := g.Model.GetPatrolCount(in.UserID)
	if err != nil {
		return nil, err
	}
	//
	return &pb.GetPatrolListByUserIDResponse{Patrol: patrols,Count: count}, nil
}

//
func (g *GrpcServer) getPoints(requestID string) ([]*structs.Point) {
	members := g.Redis.ZMembers(requestID)
	//
	zPoints := g.Redis.ZPoints(requestID, members)
	//
	var points []*structs.Point
	for _,geoPos := range zPoints {
		for _,item := range geoPos {
			points = append(points, &structs.Point{
				Longitude: item.Longitude,
				Latitude: item.Latitude,
			})
		}
	}
	return points
}

// 从redis中获取点位信息
func (g *GrpcServer) getPointsFromRedis(key string) []*structs.Point  {
	var points []*structs.Point
	// 获取 member
	zMembers := g.Redis.ZMembers(key)
	// 获取点位信息
	zPoints := g.Redis.ZPoints(key, zMembers)
	// 处理并转换为 structs.point
	for _,geoPos := range zPoints {
		for _,item := range geoPos {
			points = append(points, &structs.Point{
				Longitude: item.Longitude,
				Latitude: item.Latitude,
			})
		}
	}
	return points
}

// 获取线数据字符串
func (g *GrpcServer) getPolylineText(points []*structs.Point)  (string, error) {
	return g.Geometry.Points2Text(points, libs.Polyline)
}

// 提交巡查数据
func (g *GrpcServer) PostPatrol(_ context.Context, in *pb.PostRequest ) (*pb.PostResponse ,  error){
	//
	defer g.Redis.Conn.Close()
	//  点位信息获取
	points := g.getPoints(in.Patrol.PatrolID)
	fmt.Printf("points %#v\n",points)
	// 巡查路线信息获取
	polyline, err := g.getPolylineText(points)
	if err != nil {
		return nil, err
	}

	log.Printf("polyline %s \n", polyline)
	//
	postData := &model.Patrol{
		UserID: in.Patrol.UserID,
		Memo: in.Patrol.Memo,
		Polyline: polyline,
		RequestID: in.Patrol.PatrolID,
	}
	patrol, err := g.Model.PostPatrol(postData)
	if err != nil {
		return nil, err
	}
	// 创建时间
	createAt,_ := ptypes.TimestampProto(patrol.CreateAt)
	resp := &pb.Patrol{PatrolID: patrol.RequestID,UserID: patrol.UserID,ID: patrol.ID,CreateAt: createAt}
	return &pb.PostResponse{Patrol: resp},nil
}

//
func (g *GrpcServer)DelPatrol(_ context.Context, in *pb.DelRequest ) (*pb.DelResponse ,  error){
	return nil,nil
}

//
func (g *GrpcServer)UpdatePatrol(_ context.Context, in *pb.UpdateRequest ) (*pb.UpdateResponse ,  error){
	return nil,nil
}



// GrpcServer
func NewGrpcServer()  {
	//
	address := conf.Config.Grpc.Port //conf.Config.Grpc.Port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//
	redisAddress := fmt.Sprintf("%s:%d",conf.Config.Redis.Host,conf.Config.Redis.Port)
	redis := libs.NewRedigo(redisAddress,conf.Config.Redis.Password,conf.Config.Redis.Selected )
	log.Printf("NewGrpcServer#redisAddress#%s\n",redisAddress)
	// postgres dataSourceName
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName)
	log.Printf("NewGrpcServer#dataSourceName#%s\n",dataSourceName)
	pb.RegisterPatrolServiceServer(s, &GrpcServer{
		Model: model.NewModel(dataSourceName),
		Redis: redis,
		Geometry: libs.NewGeometry(),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Printf("#Wechat#GrpcServer#Start At# %s",address)

}


