package main

import (
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/pkg/db"
	"go_im/pkg/rpc"
	"go_im/service"
	"go_im/service/api_service"
	"go_im/service/gateway"
	"go_im/service/group_messaging"
)

func main() {
	db.Init()
	dao.Init()

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	etcd := config.Etcd.Servers

	err = gateway.InitMessageProducer(config.Nsq.Nsqd)
	if err != nil {
		panic(err)
	}
	groupManager, err := group_messaging.NewClient(&rpc.ClientOptions{
		Name:        config.GroupMessaging.Client.Name,
		EtcdServers: etcd,
	})
	if err != nil {
		panic(err)
	}
	group.SetInterfaceImpl(groupManager)
	group.SetMessageHandler(client.EnqueueMessageToDevice)

	server := api_service.NewServer(&rpc.ServerOptions{
		Name:        config.Api.Server.Name,
		Network:     config.Api.Server.Network,
		Addr:        config.Api.Server.Addr,
		Port:        config.Api.Server.Port,
		EtcdServers: etcd,
	})

	err = server.Run()

	if err != nil {
		panic(err)
	}
}
