package api

import (
	"go_im/im"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/group"
	"go_im/service/rpc"
	"math"
	"testing"
)

func TestNewServer(t *testing.T) {

	api.SetImpl(im.NewApiRouter())
	client.Manager = im.NewClientManager()
	group.Manager = im.NewGroupManager()

	op := rpc.ServerOptions{
		Network:        "tcp",
		Addr:           "localhost",
		Port:           5555,
		MaxRecvMsgSize: math.MaxInt32,
		MaxSendMsgSize: math.MaxInt32,
	}
	server := NewServer(&op)
	err := server.Run()
	panic(err)
}
