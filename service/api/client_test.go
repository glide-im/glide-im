package api

import (
	"go_im/im/message"
	"go_im/service/rpc"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&rpc.ClientOptions{
		Name:        "api",
		EtcdServers: []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
	})
	defer client.Close()
	err := client.Run()
	if err != nil {
		panic(err)
	}
	client.Handle(0, message.NewMessage(1, "api.echo", ""))
	time.Sleep(time.Second * 3)
}
