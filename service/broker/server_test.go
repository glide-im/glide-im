package broker

import (
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"go_im/service"
	"testing"
)

func RunServerForTest(port int) {

	config, _ := service.GetConfig()

	c := &rpc.ClientOptions{
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
	server := NewServer(&options, c)

	err := server.Run()
	if err != nil {
		logger.E("%v", err)
	}
}

func TestNewServer(t *testing.T) {

	RunServerForTest(9090)

}
