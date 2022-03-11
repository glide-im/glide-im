package api_service

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) Handle(ctx context.Context, r *pb_rpc.ApiHandleRequest, resp *emptypb.Empty) error {
	return api.Handle(r.GetUid(), r.GetDevice(), &message.Message{CommMessage: r.GetMessage()})
}

func (s *Server) Echo(ctx context.Context, r *pb_rpc.ApiHandleRequest, resp *pb_rpc.Response) error {
	return nil
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}
