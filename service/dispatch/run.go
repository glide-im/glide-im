package dispatch

import (
	"errors"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/service"
)

func SetupClient(configs *service.Configs) error {

	config := configs.Dispatch.Client
	options := &rpc.ClientOptions{
		Name:        config.Name,
		EtcdServers: configs.Etcd.Servers,
	}
	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	// TODO remove
	client.SetInterfaceImpl(cli)
	return nil
}

func RunServer(configs *service.Configs) error {

	sConfig := configs.Dispatch.Server
	if sConfig == nil {
		return errors.New("invalid server config")
	}

	server, err := NewServer(configs.Nsq.Nsqd, &rpc.ServerOptions{
		Name:        sConfig.Name,
		Network:     sConfig.Network,
		Addr:        sConfig.Addr,
		Port:        sConfig.Port,
		EtcdServers: configs.Etcd.Servers,
	})

	if err != nil {
		return err
	}
	return server.Run()
}
