package client

import (
	"context"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/rpc"
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

func (s *Server) ClientSignIn(ctx context.Context, request *pb.SignInRequest, reply *pb.Response) error {
	client.Manager.ClientSignIn(request.GetOld(), request.GetUid(), request.GetDevice())
	return nil
}

func (s *Server) UserLogout(ctx context.Context, request *pb.UidRequest, reply *pb.Response) error {
	client.Manager.UserLogout(request.GetUid())
	return nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.UidMessageRequest, reply *pb.Response) error {
	err := client.Manager.HandleMessage(request.GetFrom(), unwrapMessage(request.GetMessage()))
	if err != nil {
		// handle err
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
