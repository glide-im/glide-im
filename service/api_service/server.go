package api_service

import (
	"context"
	"errors"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_im"
	"go_im/protobuf/gen/pb_rpc"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) Handle(ctx context.Context, r *pb_rpc.ApiHandleRequest, resp *pb_im.CommMessage) error {
	msg := message.FromProtobuf(r.Message)

	if msg.GetAction() == "api.user.auth" {
		protobuf := msg.GetProtobuf()
		extra := protobuf.GetExtra()
		if extra == nil {
			return errors.New("message extra is nil")
		}
	}
	m, err := api.Handle(r.GetUid(), r.GetDevice(), msg)
	resp = m.GetProtobuf()
	return err
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
