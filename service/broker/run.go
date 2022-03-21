package broker

import (
	"go_im/im/group"
	"go_im/service"
)

func SetupClient(config *service.Configs) error {

	broker := config.Broker.Client
	if broker == nil {

	}

	options := config.Broker.Client.ToClientOptions()
	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	group.SetInterfaceImpl(cli)
	return nil
}

func RunServer(configs *service.Configs) error {
	options := configs.GroupMessaging.Server.ToServerOptions(configs.Etcd.Servers)
	server := NewServer(options)
	return server.Run()
}
