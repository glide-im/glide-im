package im_service

import (
	"context"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
)

type Server struct {
	srv *rpc.BaseServer
}

func RunServer(options *rpc.ServerOptions) error {
	s := &Server{
		srv: rpc.NewBaseServer(options),
	}
	s.srv.Register(options.Name, s)
	return s.srv.Run()
}

func (s *Server) ClientSignIn(ctx context.Context, request *pb_rpc.GatewaySignInRequest, reply *pb_rpc.Response) error {
	err := client.SignIn(request.GetOld(), request.GetUid(), request.GetDevice())
	if err != nil {
		reply.Message = err.Error()
		reply.Ok = false
	}
	return nil
}

func (s *Server) ClientLogout(ctx context.Context, request *pb_rpc.GatewayLogoutRequest, reply *pb_rpc.Response) error {
	err := client.Logout(request.GetUid(), request.GetDevice())
	if err != nil {
		reply.Message = err.Error()
		reply.Ok = false
	}
	return nil
}

func (s *Server) EnqueueMessage(ctx context.Context, request *pb_rpc.EnqueueMessageRequest, reply *pb_rpc.Response) error {
	m := message.FromProtobuf(request.GetMessage())
	err := client.EnqueueMessageToDevice(request.GetUid(), 0, m)
	if err != nil {
		reply.Message = err.Error()
		reply.Ok = false
	}
	return nil
}
