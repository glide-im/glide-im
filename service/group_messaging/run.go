package group_messaging

import (
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/service"
)

func SetupClient(configs *service.Configs) error {

	options := configs.GroupMessaging.Client.ToClientOptions()
	options.EtcdServers = configs.Etcd.Servers
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
