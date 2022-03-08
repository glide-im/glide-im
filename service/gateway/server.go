package gateway

import (
	"context"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/pkg/rpc"
	"go_im/protobuff/gen/pb_rpc"
)

const ServiceName = "client"

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
	client.EnqueueMessageToDevice(request.GetUid(), 0, unwrapMessage(request.Message))
	return nil
}

func unwrapMessage(pb_rpcMsg *pb_rpc.CommMessage) *message.Message {
	return &message.Message{}
}

func newResponse(ok bool, msg string) *pb_rpc.Response {
	return &pb_rpc.Response{
		Ok:      ok,
		Message: msg,
	}
}
