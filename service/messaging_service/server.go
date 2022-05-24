package messaging_service

import (
	"context"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/im/messaging"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
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

func (s *Server) HandleMessage(ctx context.Context, request *pb_rpc.MessagingHandleRequest, replay *pb_rpc.Response) error {
	m := message.FromProtobuf(request.Message)

	return messaging.HandleMessage(request.GetId(), request.GetDevice(), m)
}
