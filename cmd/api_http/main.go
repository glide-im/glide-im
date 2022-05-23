package main

import (
	"go_im/config"
	"go_im/im/api"
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/gateway"
)

func main() {

	db.Init()
	dao.Init()

	// api 需要发送 gateway 和 group 消息
	// 测试的时候 mock
	//api.MockDep()
	initIM()

	err := api.RunHttpServer("0.0.0.0", 8081)

	if err != nil {
		panic(err)
	}
}

func initIM() {

	im := config.ApiHttpService.IMService
	clientConfig := service.ClientConfig{
		Addr:        im.Addr,
		Port:        im.Port,
		EtcdServers: im.Etcd,
		Name:        im.Name,
	}
	err := gateway.SetupClient(&clientConfig)
	if err != nil {
		panic(err)
	}
}
