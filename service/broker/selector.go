package broker

import (
	"context"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
)

const ctxKeyRoute = "_key_group_route"

type groupRouteSelector struct {
	servers map[string]string

	round rpc.RoundRobinSelector
}

func newGroupRouteSelector() *groupRouteSelector {
	return &groupRouteSelector{}
}

func (s *groupRouteSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	route, ok := ctx.Value(ctxKeyRoute).(string)
	if ok {
		_, exist := s.servers[route]
		if !exist {
			logger.E("group route not exist, srv:%s", route)
			return ""
		}
		return route
	}

	return s.SelectNext()
}

func (s *groupRouteSelector) SelectNext() string {
	return s.round.SelectNext()
}

func (s *groupRouteSelector) UpdateServer(servers map[string]string) {
	s.servers = servers

	s.round.UpdateServer(servers)
}
