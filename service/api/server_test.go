package api

import (
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/group"
	"go_im/service/rpc"
	"testing"
)

func TestNewServer(t *testing.T) {

	api.SetHandler(api.NewApiRouter())
	client.Manager = client.NewClientManager()
	group.Manager = group.NewGroupManager()

	op := rpc.ServerOptions{
		Name:        "api",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8972,
		EtcdServers: []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
	}
	server := NewServer(&op)
	err := server.Run()
	t.Error(err)
}

func TestNewServer2(t *testing.T) {

	api.SetHandler(api.NewApiRouter())
	client.Manager = client.NewClientManager()
	group.Manager = group.NewGroupManager()

	op := rpc.ServerOptions{
		Name:        "api",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8973,
		EtcdServers: []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
	}
	server := NewServer(&op)
	err := server.Run()
	t.Error(err)
}
