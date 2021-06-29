package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/wechat/conf"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)
func MyInit()  {
	// 配置文件初始化
	confFile := "../conf/conf.yaml"
	err := conf.NewConf(confFile)
	if err !=nil {
		log.Fatalf("init conf unmarshal error: %s",err)
	}
}
//
type Service struct {
	ID int64
	Host string
	Name string
	Meta *Meta
}

// 服务信息
type Meta struct {
	//
	Desc string
	Contact string
	Phone string
}

// NewService
func NewServices(ID int64,host,name string, meta *Meta) ([]byte, error) {
	//
	service := &Service{
		ID: ID,
		Host: host,
		Name: name,
		Meta: meta,
	}
	return json.Marshal(service)
}
// 健康检查



var Client *clientv3.Client
//
var IDChan  chan clientv3.LeaseID
//
var ExitChan chan bool
//
func TestMain(m *testing.M)  {

	MyInit()
	fmt.Printf("conf.Config.Register.Host: %s\n",conf.Config.Register.Host)
	cli, err := clientv3.NewFromURL(conf.Config.Register.Host)
	if err !=nil {
		log.Fatalf("init etcd client error %s\n", err)
	}
	// 存储目录P
	Client = cli
	IDChan = make(chan clientv3.LeaseID)
	ExitChan = make(chan bool)
	m.Run()
}

//
func TestUnRegister(t *testing.T)  {
	//
	// 需要用一个来存储 lease.ID
	ID := <- IDChan
	t.Logf("ID %#v \n", ID)
	//
	if _ , err:= Client.Revoke(context.TODO(), ID); err !=nil {
		t.Fatalf("Error %s\n",err)
	}
	//
	ExitChan <- true
	//
	t.Logf("Revoke Success\n")

}
//
func TestRegister(t *testing.T)   {

	service := &Service{ID: 1,Host: "127.0.0.1:8081",Name:"Wechat",Meta: &Meta{Contact: "clone", Phone: "17605048999",Desc: "Wechat Service"}}
	srv, err := NewServices(service.ID,service.Host,service.Name,service.Meta)
	if err !=nil {
		t.Fatalf("error %s\n",err)
	}
	key := fmt.Sprintf("/services/%s/%s",service.Name,service.Host)
	var curLeaseId clientv3.LeaseID = 0
	lease := clientv3.NewLease(Client)
	//resp, err := lease.Grant(context.TODO(),20) // 20秒
	if err !=nil {
		t.Logf("lease error %s\n",err)
	}

	//
	//Client.Put(context.TODO(),key,string(srv))
	// ping
	go func() {

		// 发送健康检查(设置有服务的有效期 如果服务挂了则不会再更新即 20秒后 就不存在了)
		for {
			if curLeaseId == 0 {
				resp, err := lease.Grant(context.TODO(),20) // 20秒
				if err !=nil {
					log.Fatalf("lease error")
				}
				Client.Put(context.TODO(),key,string(srv),clientv3.WithLease(resp.ID) )
				curLeaseId = resp.ID
				//
				IDChan <- curLeaseId

				t.Logf("chan： %#v\n ", IDChan)
			}else {
				if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
					curLeaseId = 0
					ID := <- IDChan
					t.Logf("chan ID： %#v\n ", ID)
					continue
				}
			}

			time.Sleep(time.Second * 18)
			t.Log("time sleep 18")
		}
	}()

	//
	//service = &Service{ID: 2,Host: "127.0.0.2:8081",Name:"Wechat",Meta: &Meta{Contact: "clone", Phone: "17605048999",Desc: "Wechat Service"}}
	//srv, err = NewServices(service.ID,service.Host,service.Name,service.Meta)
	//if err !=nil {
	//	log.Fatalf("error %s\n",err)
	//}
	//key = fmt.Sprintf("/services/%s/%s",service.Name,service.Host)
	////
	//Client.Put(context.TODO(),key,string(srv),clientv3.WithLease(resp.ID) )
	//
	//service = &Service{ID: 3,Host: "127.0.0.3:8081",Name:"Wechat",Meta: &Meta{Contact: "clone", Phone: "17605048999",Desc: "Wechat Service"}}
	//srv, err = NewServices(service.ID,service.Host,service.Name,service.Meta)
	//if err !=nil {
	//	log.Fatalf("error %s\n",err)
	//}
	//key = fmt.Sprintf("/services/%s/%s",service.Name,service.Host)
	//
	//Client.Put(context.TODO(),key,string(srv),clientv3.WithLease(resp.ID) )
	// 等待
	<- ExitChan
	t.Logf("Done")
}
// 服务发现
func TestDiscover(t *testing.T)  {
	//
	key := "/services/Wechat/"
	resp, err := Client.Get(context.TODO(),key,clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("error %s\n", err)
	}
	var service Service
	var services []*Service
	//
	for _,item := range resp.Kvs {
		//
		val := item.Value
		err := json.Unmarshal([]byte(val),&service)
		if err !=nil {
			log.Fatalf("error %s\n", err)
		}
		fmt.Printf("service %#v\n", service)
		services = append(services, &service)
	}
	// 所有的服务 需要康复的服务(因在服务注册时已添加了健康检查 因此这里获取到的是健康的服务)
	fmt.Printf("%#v\n",services)
	//
}

//




func TestPut(t *testing.T)  {
	//
	fmt.Printf("HelloWorld")
	resp , err := Client.Put(context.TODO(),"/wechat/127.0.0.1:8081","127.0.0.1:8081")
	fmt.Printf("resp %#v\n error :%s\n", resp, err)
	resp , err = Client.Put(context.TODO(),"/wechat/127.0.0.2:8081","127.0.0.2:8081")
	fmt.Printf("resp %#v\n error :%s\n", resp, err)
	resp , err = Client.Put(context.TODO(),"/wechat/127.0.0.3:8081","127.0.0.3:8081")
	fmt.Printf("resp %#v\n error :%s\n", resp, err)
}

func TestGet(t *testing.T)  {
	resp, err := Client.Get(context.TODO(),"/wechat/", clientv3.WithPrefix())
	if err != nil {
		t.Logf("error %s\n", err)
	}

	for key, item := range resp.Kvs{
		t.Logf("key %d\n", key )
		t.Logf("item.key %s\n item.value %s\n",item.Key,item.Value)
	}
}




