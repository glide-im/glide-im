package group_messaging

import (
	"go_im/pkg/rpc"
	"go_im/service"
	"testing"
)

func TestNewServer(t *testing.T) {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	etcd := config.Etcd.Servers

	server := NewServer(&rpc.ServerOptions{
		Name:        config.GroupMessaging.Server.Name,
		Network:     config.GroupMessaging.Server.Network,
		Addr:        config.GroupMessaging.Server.Addr,
		Port:        config.GroupMessaging.Server.Port,
		EtcdServers: etcd,
	})

	err = server.Run()
}

func TestNewServer2(t *testing.T) {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	etcd := config.Etcd.Servers

	server := NewServer(&rpc.ServerOptions{
		Name:        config.GroupMessaging.Server.Name,
		Network:     config.GroupMessaging.Server.Network,
		Addr:        config.GroupMessaging.Server.Addr,
		Port:        config.GroupMessaging.Server.Port + 100,
		EtcdServers: etcd,
	})

	err = server.Run()
}
