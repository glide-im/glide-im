package group_messaging

import (
	"context"
	"go_im/im/message"
	"go_im/protobuff/gen/pb_rpc"
	"go_im/service/rpc"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) UpdateMember(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {
	panic("implement me")
}

func (s *Server) UpdateGroup(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {
	panic("implement me")
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {
	panic("implement me")
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {
	panic("implement me")
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}

func unwrapMessage(pbMsg *pb_rpc.CommMessage) *message.Message {
	return &message.Message{}
}
