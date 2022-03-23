package broker

import (
	"go_im/im/group"
	"go_im/pkg/logger"
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

	srvOptions := configs.Broker.Server.ToServerOptions(configs.Etcd.Servers)

	cliOptions := configs.GroupMessaging.Client.ToClientOptions()

	logger.D("broker %s", "run server")
	server, err := NewServer(srvOptions, cliOptions)
	if err != nil {
		return err
	}
	return server.Run()
}
