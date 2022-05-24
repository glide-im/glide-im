package dispatch

import (
	"github.com/glide-im/glideim/pkg/rpc"
	"testing"
)

func TestNewClient(t *testing.T) {

	options := &rpc.ClientOptions{
		Addr: "127.0.0.1",
		Port: 8080,
		Name: "dispatch",
	}
	c, err := NewClient(options)
	if err != nil {
		t.Error(err)
	}
	err = c.UpdateGatewayRoute(1, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	_ = c.Close()
}
