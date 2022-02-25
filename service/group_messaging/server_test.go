package group_messaging

import (
	"go_im/service"
	"go_im/service/rpc"
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
