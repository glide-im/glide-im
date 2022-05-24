package broker

import (
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/pkg/logger"
	"github.com/glide-im/glideim/service"
)

func SetupClient(config *service.Configs) error {

	options := config.Broker.Client.ToClientOptions()
	options.EtcdServers = config.Etcd.Servers

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
