package rpc

import (
	"context"
	"fmt"
	"github.com/glide-im/glideim/pkg/logger"
	"google.golang.org/grpc"
	"math"
	"net"
)

type Runnable interface {
	Run() error
}

type BaseGRpcServer struct {
	RpcServer *grpc.Server
	Socket    net.Listener

	AppId   int64
	Options *ServerOptions
}

func NewBaseGRpcServer(options *ServerOptions) *BaseGRpcServer {
	ret := &BaseGRpcServer{
		Options: options,
	}
	ret.init(options)
	return ret
}

func (s *BaseGRpcServer) init(options *ServerOptions) {
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
		grpc.UnaryInterceptor(s.unaryLogInterceptor),
		grpc.MaxRecvMsgSize(options.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(options.MaxSendMsgSize),
	}
	s.RpcServer = grpc.NewServer(op...)
}

func (s *BaseGRpcServer) unaryLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger.I("grpc server method called: %s", info.FullMethod)
	return handler(ctx, req)
}

func (s *BaseGRpcServer) Run() error {
	if s.Options == nil {
		s.init(nil)
	}
	return s.RpcServer.Serve(s.Socket)
}
