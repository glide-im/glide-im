package broker

import (
	"github.com/glide-im/glideim/pkg/logger"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/service"
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
