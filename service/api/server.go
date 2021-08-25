package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/api/rpc"
	rpc2 "go_im/service/rpc"
)

type Server struct {
	*rpc2.BaseServer
}

func NewServer(options *rpc2.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc2.NewBaseServer(options),
	}
	rpc.RegisterApiServiceServer(s.RpcServer, &Server{})
	return s
}

func (a *Server) Handle(ctx context.Context, request *rpc.Request) (*rpc.Response, error) {
	msg := &message.Message{
		Seq:    request.Message.Seq,
		Action: message.Action(request.Message.Action),
		Data:   request.Message.Data,
	}
	api.Handle(request.Uid, msg)
	return nil, nil
}
