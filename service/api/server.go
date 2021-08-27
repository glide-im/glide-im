package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/service/api/pb"
	"go_im/service/rpc"
)

type Server struct {
	*rpc.BaseServer
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	pb.RegisterApiServiceServer(s.RpcServer, s)
	return s
}

func (a *Server) Handle(ctx context.Context, request *pb.HandleRequest) (*pb.Response, error) {
	msg := &message.Message{
		Seq:    request.Message.Seq,
		Action: message.Action(request.Message.Action),
		Data:   request.Message.Data,
	}
	api.Handle(request.Uid, msg)
	return &pb.Response{Ok: true}, nil
}

func (a *Server) Run() error {
	logger.D("gRPC Api server run, %s@%s:%d", a.Options.Network, a.Options.Addr, a.Options.Port)
	return a.BaseServer.Run()
}
