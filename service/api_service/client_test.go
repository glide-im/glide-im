package api_service

import (
	"go_im/im/message"
	"go_im/service/rpc"
	"testing"
	"time"
)

var etcd = []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"}

func TestNewClient(t *testing.T) {
	opts := &rpc.ClientOptions{
		Name:        "api",
		EtcdServers: []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
	}
	opts.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	//opts.SerializeType = protocol.ProtoBuffer
	client, _ := NewClient(opts)
	defer client.Close()
	err := client.Run()
	if err != nil {
		panic(err)
	}
	client.Handle(0, 0, message.NewMessage(1, "api.app.echo", ""))
	time.Sleep(time.Second * 3)
}
