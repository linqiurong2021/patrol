package libs
//
import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
)

type Meta struct {
	// 版本信息
	Version string `json:"version"`
	// 作者
	Contact string `json:"contact"`
	// 联系电话
	Phone string `json:"phone"`
	// 有哪些Api
	Api []string `json:"api"`
	// 服务描述
	Desc string `json:"desc"`
}
type Service struct {
	ID int64 `json:"id"`
	Host string `json:"host"`
	Meta string `json:"meta"`
}

type Etcd struct {
	Client *clientv3.Client
	Dir string
}
//
func NewEtcd(host string, dir string) *Etcd {
	//
	cli, err := clientv3.NewFromURL(host)
	if err !=nil {
		log.Fatalf("init etcd client error %s\n", err)
	}
	// 存储目录
	if dir == "" {
		dir = "/services/"
	}
	//
	return &Etcd{
		Client:cli,
		Dir: dir,
	}
}
//
func (etcd *Etcd) Watch(key string)   {
	//
}

// 注册服务
func (etcd *Etcd) Register(key string,service *Service) (*Service, error) {
	// 判断是否存在key  如果不存在则新增  如果已存在则获取后 再添加 push
	resp , err := etcd.Client.Get(context.TODO(),key,clientv3.WithPrefix())
	if err !=nil {
		return  nil, err
	}
	var services []*Service

	// lease
	//lease := clientv3.NewLease(etcd.Client)
	//grant, err := lease.Grant(context.TODO(),5)

	// 有值需要获取后再新增
	if resp.Count >= 0 {
		//
		for _, item := range resp.Kvs {
			err = json.Unmarshal([]byte(item.Value), &services)
			if err != nil {
				log.Printf("json unmarshal err %#v\n", err)
				return nil, err
			}
		}
		//
		services = append(services,service)
		//
		bytes, err := json.Marshal(&services)
		if err !=nil {
			log.Printf("json marshal err %#v\n", err)
			return nil, err
		}
		// 更新数据
		etcd.Client.Put(context.TODO(), key,string(bytes))
	}else{
		services = append(services, service)
		bytes, err := json.Marshal(services)
		if err != nil {
			log.Printf("json marshal err %s\n", err)
		}
		//
		resp ,err := etcd.Client.Put(context.TODO(),key,string(bytes))
		fmt.Sprintf("new service resp %#v\n",resp)
	}

	return service, nil
}

// 获取一个服务
func (etcd *Etcd) GetServer(key string)  (*Service, error)  {

	return nil, nil
}

// 发现服务
func (etcd *Etcd) Discover(key string)  (*Service, error)  {
	//
	return nil, nil
}

// 负载均衡
func (etcd *Etcd) Balance(key string) (*Service, error)  {
	return nil, nil
}




