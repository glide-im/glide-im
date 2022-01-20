package group

import (
	"context"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/rpc"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) UpdateMember(ctx context.Context, request interface{}, replay interface{}) error {
	panic("implement me")
}

func (s *Server) UpdateGroup(ctx context.Context, request interface{}, replay interface{}) error {
	panic("implement me")
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request interface{}, replay interface{}) error {
	panic("implement me")
}

func (s *Server) DispatchMessage(ctx context.Context, request interface{}, replay interface{}) error {
	panic("implement me")
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}

func unwrapMessage(pbMsg *pb.Message) *message.Message {
	return &message.Message{
		Seq:    pbMsg.Seq,
		Action: message.Action(pbMsg.Action),
	}
}
