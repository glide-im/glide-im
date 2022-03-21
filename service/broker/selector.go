package broker

import "context"

var routeCache map[int64]string

type groupRouteSelector struct {
}

func (s *groupRouteSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	return ""
}

func (s *groupRouteSelector) UpdateServer(servers map[string]string) {

}
