package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
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
