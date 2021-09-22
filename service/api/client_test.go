package api

import (
	"github.com/stretchr/testify/assert"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/service/route"
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
	client := NewClient(opts)
	defer client.Close()
	err := client.Run()
	if err != nil {
		panic(err)
	}
	client.Handle(0, message.NewMessage(1, "api.app.echo", ""))
	time.Sleep(time.Second * 3)
}

func TestRegisterToRoute(t *testing.T) {
	assert.Nil(t, route.RegisterService("api", etcd))
}

func TestNewClientByRouter(t *testing.T) {
	cli := NewClientByRouter("api", &rpc.ClientOptions{
		Name:        "route",
		EtcdServers: etcd,
	})
	defer cli.Close()

	logger.D("=%s", "1")
	for i := 0; i < 3; i++ {
		cli.Handle(1, &message.Message{
			Seq:    1,
			Action: "api.app.echo",
			Data:   "this is data",
		})
	}
}
