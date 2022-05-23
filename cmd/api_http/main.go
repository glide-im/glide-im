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

	err := config.Load()
	if err != nil {
		panic(err)
	}

	db.Init()
	dao.Init()

	// api 需要发送 gateway 和 group 消息
	// 测试的时候 mock
	//api.MockDep()
	initIM()

	addr := config.ApiHttp.Addr
	port := config.ApiHttp.Port
	err = api.RunHttpServer(addr, port)

	if err != nil {
		panic(err)
	}
}

func initIM() {

	addr := config.IMRpcServer.Addr
	port := config.IMRpcServer.Port
	clientConfig := service.ClientConfig{
		Addr: addr,
		Port: port,
	}
	err := gateway.SetupClient(&clientConfig)
	if err != nil {
		panic(err)
	}
}
