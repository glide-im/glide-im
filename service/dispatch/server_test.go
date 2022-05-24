package dispatch

import (
	"github.com/glide-im/glideim/pkg/rpc"
	"testing"
)

func TestNewServer(t *testing.T) {

	options := &rpc.ServerOptions{
		Name:    "dispatch",
		Network: "tcp",
		Addr:    "0.0.0.0",
		Port:    8080,
	}
	server, err := NewServer("127.0.0.1:4154", options)
	if err != nil {
		t.Error(err)
		return
	}
	err = server.Run()
	t.Error(err)
}
