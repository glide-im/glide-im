package client

import (
	"context"
	"go_im/im"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/pkg/logger"
	pb2 "go_im/service/api/pb"
	"go_im/service/client/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	*rpc.BaseServer
	im im.Server
}

var apiServiceClient pb2.ApiServiceClient

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	pb.RegisterClientServiceServer(s.RpcServer, s)
	return s
}

func (s *Server) ClientSignIn(ctx context.Context, request *pb.SignInRequest) (*pb.Response, error) {
	client.Manager.ClientSignIn(request.GetOld(), request.GetUid(), request.GetDevice())
	return newResponse(true, "ok"), nil
}

func (s *Server) UserLogout(ctx context.Context, request *pb.UidRequest) (*pb.Response, error) {
	client.Manager.UserLogout(request.GetUid())
	return newResponse(true, "ok"), nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.UidMessageRequest) (*pb.Response, error) {
	err := client.Manager.DispatchMessage(request.GetFrom(), unwrapMessage(request.GetMessage()))
	if err != nil {
		// handle err
	}
	return newResponse(true, "ok"), nil
}

func (s *Server) Api(ctx context.Context, request *pb.UidMessageRequest) (*pb.Response, error) {
	_, _ = apiServiceClient.Handle(ctx, &pb2.HandleRequest{
		Uid: request.From,
		Message: &pb2.Message{
			Seq:    request.GetMessage().GetSeq(),
			Action: request.GetMessage().GetAction(),
			Data:   request.GetMessage().GetData(),
		},
	})
	return newResponse(true, "ok"), nil
}

func (s *Server) EnqueueMessage(ctx context.Context, request *pb.UidMessageRequest) (*pb.Response, error) {
	client.EnqueueMessage(request.GetFrom(), unwrapMessage(request.Message))
	return newResponse(true, "ok"), nil
}

func (s *Server) IsOnline(ctx context.Context, request *pb.UidRequest) (*pb.Response, error) {
	_ = client.Manager.IsOnline(request.GetUid())

	return newResponse(true, "ok"), nil
}

func (s *Server) Update(ctx context.Context, empty *emptypb.Empty) (*pb.Response, error) {
	client.Manager.Update()
	return newResponse(true, "ok"), nil
}

func (s *Server) Run() error {
	logger.D("gRPC Client server run, %s@%s:%d", s.Options.Network, s.Options.Addr, s.Options.Port)
	return s.BaseServer.Run()
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
