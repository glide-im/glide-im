package im_service

import (
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"testing"
)

func TestNewClient(t *testing.T) {

	options := rpc.ClientOptions{
		Addr: "0.0.0.0",
		Port: 9081,
		Name: "im_service",
	}
	cli, err := NewClient(&options)
	defer cli.Close()
	if err != nil {
		t.Error(err)
	}

	err = cli.EnqueueMessage(1, 1, message.NewEmptyMessage())

	if err != nil {
		t.Error(err)
	}
}
