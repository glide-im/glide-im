package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) Handle(ctx context.Context, r *pb.HandleRequest, resp *emptypb.Empty) error {

	api.Handle(r.Uid, &message.Message{
		Seq:    r.GetMessage().GetSeq(),
		Action: message.Action(r.GetMessage().GetAction()),
		Data:   r.GetMessage().GetData(),
	})
	return nil
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}
