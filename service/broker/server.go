package broker

import (
	"context"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
	"go_im/service/group_messaging"
)

type Server struct {
	group_messaging.Server

	selector *groupRouteSelector
	cli      *group_messaging.Client
}

func NewServer(options *rpc.ServerOptions, groupMessagingOpts *rpc.ClientOptions) *rpc.BaseServer {
	s := rpc.NewBaseServer(options)
	brokerServer := &Server{}

	brokerServer.selector = &groupRouteSelector{}
	groupMessagingOpts.Selector = brokerServer.selector

	brokerServer.cli, _ = group_messaging.NewClient(groupMessagingOpts)
	s.Register(options.Name, brokerServer)
	return s
}

func (s *Server) UpdateMember(ctx context.Context, param *pb_rpc.UpdateMemberParam, replay *pb_rpc.Response) error {
	return s.cli.Call(ctx, "UpdateMember", param, replay)
}

func (s *Server) UpdateGroup(ctx context.Context, param *pb_rpc.UpdateGroupParam, replay *pb_rpc.Response) error {

	return s.cli.Call(ctx, "UpdateGroup", param, replay)
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, param *pb_rpc.DispatchGroupNotifyParam, replay *pb_rpc.Response) error {

	return s.cli.Call(ctx, "DispatchNotifyMessage", param, replay)
}

func (s *Server) DispatchMessage(ctx context.Context, param *pb_rpc.DispatchGroupChatParam, replay *pb_rpc.Response) error {

	return s.cli.Call(ctx, "DispatchMessage", param, replay)
}
