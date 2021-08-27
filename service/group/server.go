package group

import (
	"context"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/service/group/pb"
	"go_im/service/rpc"
)

type Server struct {
	*rpc.BaseServer
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	pb.RegisterGroupServiceServer(s.RpcServer, s)
	return s
}

func (s *Server) PutMember(ctx context.Context, request *pb.PutMemberRequest) (*pb.Response, error) {
	gm := pbMember2daoMember(request.GetMember())[0]
	group.Manager.PutMember(request.GetGid(), gm)
	return newResponse(true, "ok"), nil
}

func (s *Server) RemoveMember(ctx context.Context, request *pb.RemoveMemberRequest) (*pb.Response, error) {
	err := group.Manager.RemoveMember(request.Gid, request.Uid...)
	if err != nil {

	}
	return newResponse(true, ""), nil
}

func (s *Server) GetMembers(ctx context.Context, request *pb.GidRequest) (*pb.GetMembersResponse, error) {

	members, err := group.Manager.GetMembers(request.GetGid())
	if err != nil {

	}
	ret := pb.GetMembersResponse{Members: daoMember2pbMember(members...)}

	return &ret, nil
}

func (s *Server) AddGroup(ctx context.Context, request *pb.AddGroupRequest) (*pb.Response, error) {

	g := pbGroup2daoGroup(request.GetGroup())
	owner := pbMember2daoMember(request.GetOwner())[0]
	group.Manager.AddGroup(g, request.GetCid(), owner)
	return newResponse(true, ""), nil
}

func (s *Server) GetGroup(ctx context.Context, request *pb.GidRequest) (*pb.Group, error) {
	g := group.Manager.GetGroup(request.GetGid())
	return daoGroup2pbGroup(g), nil
}

func (s *Server) GetGroupCid(ctx context.Context, request *pb.GidRequest) (*pb.GetCidResponse, error) {
	cid := group.Manager.GetGroupCid(request.GetGid())
	return &pb.GetCidResponse{Cid: cid}, nil
}

func (s *Server) HasMember(ctx context.Context, request *pb.HasMemberRequest) (*pb.HasMemberResponse, error) {
	return &pb.HasMemberResponse{Has: group.Manager.HasMember(request.GetGid(), request.GetUid())}, nil
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, request *pb.NotifyRequest) (*pb.Response, error) {
	group.Manager.DispatchNotifyMessage(request.Uid, request.GetGid(), unwrapMessage(request.GetMessage()))
	return newResponse(true, ""), nil
}

func (s *Server) DispatchMessage(ctx context.Context, request *pb.DispatchMessageRequest) (*pb.Response, error) {
	err := group.Manager.DispatchMessage(request.Uid, unwrapMessage(request.GetMessage()))
	if err != nil {

	}
	return newResponse(true, ""), nil
}

func (s *Server) Run() error {
	logger.D("gRPC Group server run, %s@%s:%d", s.Options.Network, s.Options.Addr, s.Options.Port)
	return s.BaseServer.Run()
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
		Members:  daoMember2pbMember(g.Members...),
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
		Members:  pbMember2daoMember(g.Members...),
	}
}

func daoMember2pbMember(members ...*dao.GroupMember) []*pb.GroupMember {
	var gm []*pb.GroupMember
	for _, member := range members {
		gm = append(gm, &pb.GroupMember{
			Id:     member.Id,
			Gid:    member.Gid,
			Uid:    member.Uid,
			Mute:   member.Mute,
			Type:   int32(member.Type),
			Remark: member.Remark,
		})
	}
	return gm
}

func pbMember2daoMember(members ...*pb.GroupMember) []*dao.GroupMember {
	var gm []*dao.GroupMember
	for _, member := range members {
		gm = append(gm, &dao.GroupMember{
			Id:     member.Id,
			Gid:    member.Gid,
			Uid:    member.Uid,
			Mute:   member.Mute,
			Type:   int8(member.Type),
			Remark: member.Remark,
		})
	}
	return gm
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
