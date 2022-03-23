package messaging_service

import (
	"go_im/im/client"
	"go_im/im/messaging"
	"go_im/service"
)

func SetupClient(configs *service.Configs) error {

	options := configs.Messaging.Client.ToClientOptions()
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
