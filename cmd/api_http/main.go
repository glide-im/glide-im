package main

import (
	"github.com/glide-im/glideim/config"
	"github.com/glide-im/glideim/im/api"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/dao"
	"github.com/glide-im/glideim/pkg/db"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/service/im_service"
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
