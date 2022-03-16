package gateway

import (
	"go_im/im/client"
	"go_im/service"
)

func SetupClient(configs service.Configs) error {

	options := configs.Gateway.Client.ToClientOptions()
	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	client.SetInterfaceImpl(cli)
	return nil
}

func RunServer(configs *service.Configs) error {

	options := configs.Gateway.Server.ToServerOptions(configs.Etcd.Servers)
	server := NewServer(options)
	return server.Run()
}
