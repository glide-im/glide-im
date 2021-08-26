package rpc

import (
	"fmt"
	"google.golang.org/grpc"
	"math"
	"net"
)

type ServerOptions struct {
	Network        string
	Addr           string
	Port           int
	MaxRecvMsgSize int
	MaxSendMsgSize int
}

type BaseServer struct {
	RpcServer *grpc.Server
	Socket    net.Listener

	options *ServerOptions
}

func NewBaseServer(options *ServerOptions) *BaseServer {
	ret := &BaseServer{
		options: options,
	}
	ret.init(options)
	return ret
}

func (s *BaseServer) init(options *ServerOptions) {
	if options == nil {
		options = &ServerOptions{
			Network:        "tcp",
			Addr:           "localhost",
			Port:           5555,
			MaxRecvMsgSize: math.MaxInt32,
			MaxSendMsgSize: math.MaxInt32,
		}
	}

	var err error
	addr := fmt.Sprintf("%s:%d", options.Addr, options.Port)
	s.Socket, err = net.Listen(options.Network, addr)
	if err != nil {
		panic(err)
	}
	op := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(options.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(options.MaxSendMsgSize),
	}
	s.RpcServer = grpc.NewServer(op...)
}

func (s *BaseServer) Run() error {
	if s.options == nil {
		s.init(nil)
	}
	return s.RpcServer.Serve(s.Socket)
}
