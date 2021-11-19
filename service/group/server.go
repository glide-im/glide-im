package group

import (
	"context"
	"go_im/im/group"
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

func (s *Server) PutMember(ctx context.Context, request *pb.PutMemberRequest, reply *pb.Response) error {
	//group.Manager.PutMember(request.GetGid(), request.GetMember())
	return nil
}

func (s *Server) RemoveMember(ctx context.Context, request *pb.RemoveMemberRequest, reply *pb.Response) error {
	//err := group.Manager.RemoveMember(request.Gid, request.Uid...)
	//if err != nil {
	//
	//}
	return nil
}

func (s *Server) ChangeStatus(ctx context.Context, request *pb.GroupStateRequest, reply *pb.Response) error {
	//group.Manager.ChangeStatus(request.GetGid(), request.GetStatus())
	return nil
}

func (s *Server) RemoveGroup(ctx context.Context, request *pb.GroupIDRequest, reply *pb.Response) error {
	//group.Manager.RemoveGroup(request.GetGid())
	return nil
}

func (s *Server) AddGroup(ctx context.Context, request *pb.GroupIDRequest, reply *pb.Response) error {
	//group.Manager.AddGroup(request.GetGid())
	return nil
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request *pb.NotifyRequest, reply *pb.Response) error {
	group.Manager.DispatchNotifyMessage(request.GetGid(), unwrapMessage(request.GetMessage()))
	return nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.DispatchMessageRequest, reply *pb.Response) error {
	//msg := request.GetMessage()
	//group.Manager.DispatchMessage(request.GetGid(), &message.GroupMessage{
	//	TargetId:    msg.GetTargetId(),
	//	Sender:      msg.GetSender(),
	//	Cid:         msg.GetCid(),
	//	UcId:        msg.GetUcId(),
	//	MessageType: int8(msg.GetMessageType()),
	//	Message:     msg.GetMessage(),
	//	SendAt:      msg.GetSendAt(),
	//})
	return nil
}

func unwrapMessage(pbMsg *pb.Message) *message.Message {
	return &message.Message{
		Seq:    pbMsg.Seq,
		Action: message.Action(pbMsg.Action),
		Data:   pbMsg.Data,
	}
}
