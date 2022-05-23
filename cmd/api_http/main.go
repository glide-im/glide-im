package main

import (
	"go_im/config"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/pkg/rpc"
	"go_im/service/im_service"
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

	options := rpc.ClientOptions{
		Addr: config.IMRpcServer.Addr,
		Port: config.IMRpcServer.Port,
		Name: config.IMRpcServer.Name,
	}
	cli, err := im_service.NewClient(&options)
	if err != nil {
		panic(err)
	}
	client.SetInterfaceImpl(cli)
}
