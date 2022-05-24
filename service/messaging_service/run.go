package messaging_service

import (
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/messaging"
	"github.com/glide-im/glideim/service"
)

func SetupClient(configs *service.Configs) error {

	options := configs.Messaging.Client.ToClientOptions()
	options.EtcdServers = configs.Etcd.Servers

	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	messaging.SetInterfaceImpl(cli.HandleMessage)
	client.SetMessageHandler(cli.HandleMessage)
	return nil
}

func RunServer(configs *service.Configs) error {

	options := configs.Messaging.Server.ToServerOptions(configs.Etcd.Servers)

	server := NewServer(options)

	return server.Run()
}
