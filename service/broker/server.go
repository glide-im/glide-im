package broker

import (
	"context"
	"go_im/im/group"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
	"go_im/service/group_messaging"
)

const ctxKeyGid = "_key_gid"

type Server struct {
	selector *groupRouteSelector
	cli      *group_messaging.Client

	routeCache map[int64]string
}

func NewServer(options *rpc.ServerOptions, groupMessagingOpts *rpc.ClientOptions) *rpc.BaseServer {
	s := rpc.NewBaseServer(options)
	brokerServer := &Server{}

	brokerServer.selector = newGroupRouteSelector()
	groupMessagingOpts.Selector = brokerServer.selector

	brokerServer.cli, _ = group_messaging.NewClient(groupMessagingOpts)
	s.Register(options.Name, brokerServer)
	return s
}

func (s *Server) UpdateMember(ctx context.Context, param *pb_rpc.UpdateMemberParam, replay *pb_rpc.Response) error {
	return s.call("UpdateMember", param, replay)
}

func (s *Server) UpdateGroup(ctx context.Context, param *pb_rpc.UpdateGroupParam, replay *pb_rpc.Response) error {

	if param.GetFlag() == group.FlagGroupCreate {
		next := s.selector.SelectNext()
		s.routeCache[param.GetGid()] = next
	}

	if param.GetFlag() == group.FlagGroupDissolve {
		delete(s.routeCache, param.GetGid())
	}

	return s.call("UpdateGroup", param, replay)
}

func (s *Server) DispatchNotifyMessage(ctx context.Context, param *pb_rpc.DispatchGroupNotifyParam, replay *pb_rpc.Response) error {
	return s.call("DispatchNotifyMessage", param, replay)
}

func (s *Server) DispatchMessage(ctx context.Context, param *pb_rpc.DispatchGroupChatParam, replay *pb_rpc.Response) error {
	return s.call("DispatchMessage", param, replay)
}

func (s *Server) call(fn string, param interface{}, replay interface{}) error {

	g, ok := param.(gidParam)
	ctx := context.TODO()
	if ok {
		gid := g.GetGid()
		route, ok := s.routeCache[gid]
		if ok {
			ctx = context.WithValue(ctx, ctxKeyRoute, route)
		} else {
			logger.E("can not find route for gid %d", gid)
		}
	}

	return s.cli.Call(ctx, fn, param, replay)
}
