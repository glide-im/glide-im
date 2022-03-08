package group_messaging

import (
	"context"
	"go_im/im/message"
	"go_im/pkg/rpc"
	"go_im/protobuff/gen/pb_rpc"
)

type Server struct {
	*rpc.BaseServer
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}

func (s *Server) UpdateMember(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {

	return nil
}

func (s *Server) UpdateGroup(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {

	return nil
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {

	return nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb_rpc.CommMessage, replay *pb_rpc.Response) error {

	return nil
}

func unwrapMessage(pbMsg *pb_rpc.CommMessage) *message.Message {
	return &message.Message{}
}
