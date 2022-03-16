package group_messaging

import (
	"go_im/im/client"
	"go_im/im/group"
	"go_im/service"
)

func SetupClient(configs *service.Configs) error {

	options := configs.GroupMessaging.Client.ToClientOptions()
	groupManager, err := NewClient(options)
	if err != nil {
		return err
	}
	group.SetInterfaceImpl(groupManager)
	group.SetMessageHandler(client.EnqueueMessageToDevice)
	return nil
}

func RunServer(configs *service.Configs) error {
	options := configs.GroupMessaging.Server.ToServerOptions(configs.Etcd.Servers)
	server := NewServer(options)
	return server.Run()
}
