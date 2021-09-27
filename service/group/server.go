package group

import (
	"context"
	"go_im/im/dao"
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
	group.Manager.PutMember(request.GetGid(), request.GetMember())
	return nil
}

func (s *Server) RemoveMember(ctx context.Context, request *pb.RemoveMemberRequest, reply *pb.Response) error {
	err := group.Manager.RemoveMember(request.Gid, request.Uid...)
	if err != nil {

	}
	return nil
}

func (s *Server) AddGroup(ctx context.Context, request *pb.AddGroupRequest, reply *pb.Response) error {

	g := pbGroup2daoGroup(request.GetGroup())
	group.Manager.AddGroup(g, request.GetOwner())
	return nil
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request *pb.NotifyRequest, reply *pb.Response) error {
	group.Manager.DispatchNotifyMessage(request.Uid, request.GetGid(), unwrapMessage(request.GetMessage()))
	return nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.DispatchMessageRequest, reply *pb.Response) error {
	err := group.Manager.DispatchMessage(request.Uid, unwrapMessage(request.GetMessage()))
	if err != nil {

	}
	return nil
}

func daoGroup2pbGroup(g *dao.Group) *pb.Group {
	return &pb.Group{
		Gid:      g.Gid,
		Name:     g.Name,
		Avatar:   g.Avatar,
		Owner:    g.Owner,
		Mute:     g.Mute,
		Notice:   g.Notice,
		CreateAt: 0,
	}
}

func pbGroup2daoGroup(g *pb.Group) *dao.Group {
	return &dao.Group{
		Gid:      g.GetGid(),
		Name:     g.GetName(),
		Avatar:   g.GetAvatar(),
		Owner:    g.GetOwner(),
		Mute:     g.GetMute(),
		Notice:   g.GetNotice(),
		CreateAt: dao.Timestamp{},
	}
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
