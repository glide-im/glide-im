package main

import (
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/pkg/db"
	"go_im/pkg/rpc"
	"go_im/service"
	"go_im/service/gateway"
	"go_im/service/group_messaging"
	"go_im/service/messaging_service"
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

	server := messaging_service.NewServer(&rpc.ServerOptions{
		Name:        config.MessageRouter.Server.Name,
		Network:     config.MessageRouter.Server.Network,
		Addr:        config.MessageRouter.Server.Addr,
		Port:        config.MessageRouter.Server.Port,
		EtcdServers: etcd,
	})
	err = server.Run()

	if err != nil {
		panic(err)
	}
}
