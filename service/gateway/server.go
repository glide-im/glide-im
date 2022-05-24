package gateway

import (
	"context"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
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

func (s *Server) ClientSignIn(ctx context.Context, request *pb_rpc.GatewaySignInRequest, reply *pb_rpc.Response) error {
	return client.SignIn(request.GetOld(), request.GetUid(), request.GetDevice())
}

func (s *Server) ClientLogout(ctx context.Context, request *pb_rpc.GatewayLogoutRequest, reply *pb_rpc.Response) error {
	return client.Logout(request.GetUid(), request.GetDevice())
}

func (s *Server) EnqueueMessage(ctx context.Context, request *pb_rpc.EnqueueMessageRequest, reply *pb_rpc.Response) error {
	m := message.FromProtobuf(request.GetMessage())
	err := client.EnqueueMessageToDevice(request.GetUid(), 0, m)
	return err
}
