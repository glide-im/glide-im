package broker

import (
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/service"
	"testing"
)

func TestNewClient(t *testing.T) {

	config, _ := service.GetConfig()

	options := rpc.ClientOptions{
		Name: config.Broker.Server.Name,
		Addr: config.Broker.Server.Addr,
		Port: config.Broker.Server.Port,
	}

	client, err := NewClient(&options)
	if err != nil {
		t.Error(err)
	}

	err = client.UpdateGroup(1, group.Update{
		Flag: group.FlagGroupCreate,
	})

	if err != nil {
		t.Error(err)
	}

	chatMessage := message.NewChatMessage(1, 1, 1, 1, 1, "", 1)
	err = client.DispatchMessage(1, "", &chatMessage)

	if err != nil {
		t.Error(err)
	}

	_ = client.Close()
}
