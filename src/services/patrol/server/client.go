package server

import (
	"context"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/golang/protobuf/ptypes"
	"github.com/linqiurong2021/patrol/src/services/patrol/conf"
	"github.com/linqiurong2021/patrol/src/services/patrol/libs"
	pb "github.com/linqiurong2021/patrol/src/services/patrol/protobuf/v1"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
	"google.golang.org/grpc"
	"log"
)

type GrpcClient struct {
	Client pb.PatrolServiceClient
	Connect *grpc.ClientConn
	Geometry *libs.Geometry
}

func NewGrpcClient() *GrpcClient{
	// 端口号
	address := conf.Config.Grpc.Port
	// 拨号
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewPatrolServiceClient(conn)
	return  &GrpcClient{
		Client: client,
		Connect: conn,
		Geometry: libs.NewGeometry(),
	}
}

//
func (c *GrpcClient) GetPatrol(ctx context.Context ,ID int64)  (*structs.Patrol, error) {

	// 请求数据
	req := pb.GetPatrolByIDRequest{
		ID: ID,
	}
	log.Printf("GetPatrol req params %#v\n",&req)
	//
	resp ,err := c.Client.GetPatrol(ctx,&req)
	//
	fmt.Sprintf("%#v,#########",resp)
	if err != nil {
		log.Printf("client login error %s\n", err)
		return nil ,err
	}
	if resp.Patrol.ID == 0 {
		return nil, xorm.ErrNotExist
	}
	line := resp.Patrol.Line
	//
	points := c.Geometry.Polyline2Points(line)
	//
	createAt, _ :=  ptypes.Timestamp(resp.Patrol.CreateAt)
	updateAt, _ := ptypes.Timestamp(resp.Patrol.UpdateAt)

	return &structs.Patrol{
		ID: resp.Patrol.ID,
		Line: points,
		CreateAt: createAt ,
		UpdateAt:updateAt,
	},nil
}

//
func (c *GrpcClient) PostPatrol(ctx context.Context ,in *structs.Patrol)  (*structs.Patrol, error) {

	// 请求数据
	req := pb.PostRequest{
		Patrol: &pb.Patrol{
			PatrolID: in.RequestID,
			Memo: in.Memo,
			UserID: 1,// 需要通过token获取
		},
	}
	log.Printf("PostPatrol req params %#v\n",req)
	//
	resp ,err := c.Client.PostPatrol(ctx,&req)
	//
	fmt.Sprintf("%#v,#########",resp)
	if err != nil {
		log.Printf("PostPatrol error %s\n", err)
		return nil ,err
	}

	return nil,nil
}



// 获取列表分页
func (c *GrpcClient) GetPatrolList(ctx context.Context , List *structs.GetPatrolListRequest) (*structs.GetPatrolListResponse, error)  {
	// 请求数据
	req := pb.GetPatrolListByUserIDRequest{UserID: List.UserID,Page: List.Page,PageSize: List.PageSize}
	//
	resp ,err := c.Client.GetPatrolList(ctx,&req)
	//
	fmt.Sprintf("%#v,#########",resp)
	if err != nil {
		fmt.Printf("client code2session error %s\n", err)
		//
	}
	fmt.Printf("共%d行\n",len(resp.Patrol))
	//
	var patrol []*structs.Patrol
	//
	for _,item := range resp.Patrol {
		createAt, _ :=  ptypes.Timestamp(item.CreateAt)
		updateAt, _ := ptypes.Timestamp(item.UpdateAt)
		//
		//var points []structs.Point
		// 列表就不做转换了
		// 点位信息转换
		patrol = append(patrol, &structs.Patrol{
			ID: item.ID,
			CreateAt: createAt,
			UpdateAt: updateAt,
		})
	}

	return &structs.GetPatrolListResponse{
		List: patrol,
		Count: resp.Count,
	}, err
}

