package gateway_service

import (
	"context"
	"fmt"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/protobuff/pb_rpc"
	"go_im/service/rpc"
)

const ServiceName = "client"

type Server struct {
	*rpc.BaseServer
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	var err error
	myAddr := fmt.Sprintf("%s@%s:%d", options.Network, options.Addr, options.Port)
	client.Manager, err = newManager(options.EtcdServers, myAddr)
	if err != nil {
		return nil
	}
	s.Register(options.Name, s)
	return s
}

func (s *Server) ClientSignIn(ctx context.Context, request *pb_rpc.GatewaySignInRequest, reply *pb_rpc.Response) error {
	client.Manager.ClientSignIn(request.GetOld(), request.GetUid(), request.GetDevice())
	return nil
}

func (s *Server) ClientLogout(ctx context.Context, request *pb_rpc.GatewayLogoutRequest, reply *pb_rpc.Response) error {
	client.Manager.ClientLogout(request.GetUid(), request.GetDevice())
	return nil
}

func (s *Server) EnqueueMessage(ctx context.Context, request *pb_rpc.EnqueueMessageRequest, reply *pb_rpc.Response) error {
	client.Manager.EnqueueMessage(request.GetUid(), 0, unwrapMessage(request.Message))
	return nil
}

func unwrapMessage(pb_rpcMsg *pb_rpc.Message) *message.Message {
	return &message.Message{}
}

func newResponse(ok bool, msg string) *pb_rpc.Response {
	return &pb_rpc.Response{
		Ok:      ok,
		Message: msg,
	}
}
