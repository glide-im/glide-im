package group_messaging

import (
	"context"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
)

type Server struct {
}

func NewServer(options *rpc.ServerOptions) *rpc.BaseServer {
	s := rpc.NewBaseServer(options)
	s.Register(options.Name, &Server{})
	return s
}

func (s *Server) UpdateMember(ctx context.Context, param *pb_rpc.UpdateMemberParam, replay *pb_rpc.Response) error {
	gid := param.GetGid()
	var updates = make([]group.MemberUpdate, len(param.GetUpdates()))
	for _, u := range param.GetUpdates() {
		updates = append(updates, group.MemberUpdate{
			Uid:   u.GetUid(),
			Flag:  u.GetFlag(),
			Extra: nil,
		})
	}
	return group.UpdateMember(gid, updates)
}

func (s *Server) UpdateGroup(ctx context.Context, param *pb_rpc.UpdateGroupParam, replay *pb_rpc.Response) error {
	gid := param.GetGid()
	update := group.Update{
		Flag:  param.GetFlag(),
		Extra: nil,
	}
	err := group.UpdateGroup(gid, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, param *pb_rpc.DispatchGroupNotifyParam, replay *pb_rpc.Response) error {
	gid := param.GetGid()
	n := param.GetNotify()
	notify := message.GroupNotify{GroupNotify: n}
	return group.DispatchNotifyMessage(gid, &notify)
}

func (s *Server) DispatchMessage(ctx context.Context, param *pb_rpc.DispatchGroupChatParam, replay *pb_rpc.Response) error {
	gid := param.GetGid()
	m := param.GetMessage()
	chatMessage := message.ChatMessage{ChatMessage: m}
	logger.D("%v,%v", gid, chatMessage)
	return nil
	//return group.DispatchMessage(gid, &chatMessage)
}
