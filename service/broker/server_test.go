package broker

import (
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"go_im/service"
	"testing"
)

func TestNewServer(t *testing.T) {

	config, _ := service.GetConfig()

	cli := &rpc.ClientOptions{
		Addr: config.GroupMessaging.Server.Addr,
		Port: config.GroupMessaging.Server.Port,
		Name: config.GroupMessaging.Server.Name,
	}

	broker := config.Broker.Server
	options := rpc.ServerOptions{
		Name:    broker.Name,
		Network: broker.Network,
		Addr:    broker.Addr,
		Port:    broker.Port,
	}

	t.Logf("Starting server on %s:%d", options.Addr, options.Port)
	server, err := NewServer(&options, cli)
	if err != nil {
		t.Error(err)
	}

	t.Log("Starting server...")
	err = server.Run()
	if err != nil {
		logger.E("%v", err)
	}

}
