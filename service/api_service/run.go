package api_service

import (
	"go_im/im/api"
	"go_im/service"
)

func SetupClient(configs *service.Configs) error {

	options := configs.Api.Client.ToClientOptions()
	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	api.SetInterfaceImpl(cli)
	return nil
}

func RunServer(configs *service.Configs) error {

	options := configs.Api.Server.ToServerOptions(configs.Etcd.Servers)
	server := NewServer(options)
	return server.Run()
}
