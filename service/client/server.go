package client

import (
	"context"
	"fmt"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/service/pb"
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

func (s *Server) ClientSignIn(ctx context.Context, request *pb.SignInRequest, reply *pb.Response) error {
	client.Manager.ClientSignIn(request.GetOld(), request.GetUid(), request.GetDevice())
	return nil
}

func (s *Server) ClientLogout(ctx context.Context, request *pb.UidRequest, reply *pb.Response) error {
	client.Manager.ClientLogout(request.GetUid())
	return nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.UidMessageRequest, reply *pb.Response) error {
	err := client.Manager.HandleMessage(request.GetFrom(), unwrapMessage(request.GetMessage()))
	if err != nil {
		// handle err
		return err
	}
	return nil
}

func (s *Server) EnqueueMessage(ctx context.Context, request *pb.UidMessageRequest, reply *pb.Response) error {
	client.EnqueueMessage(request.GetFrom(), unwrapMessage(request.Message))
	return nil
}

func unwrapMessage(pbMsg *pb.Message) *message.Message {
	return &message.Message{
		Seq:    pbMsg.Seq,
		Action: message.Action(pbMsg.Action),
		Data:   pbMsg.Data,
	}
}

func newResponse(ok bool, msg string) *pb.Response {
	return &pb.Response{
		Ok:      ok,
		Message: msg,
	}
}
