package api

import (
	"go_im/im/message"
	"go_im/service/rpc"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	opts := &rpc.ClientOptions{
		Name:        "api",
		EtcdServers: []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
	}
	opts.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	//opts.SerializeType = protocol.ProtoBuffer
	client := NewClient(opts)
	defer client.Close()
	err := client.Run()
	if err != nil {
		panic(err)
	}
	client.Handle(0, message.NewMessage(1, "api.app.echo", ""))
	time.Sleep(time.Second * 3)
}
