package route

import (
	"context"
	"github.com/smallnest/rpcx/client"
)

type router struct {
	services map[string]string
}

func (r *router) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	for k := range r.services {
		println(k)
		return k
	}
	return ""
}

func (r *router) UpdateServer(servers map[string]string) {
	for k, v := range servers {
		r.services[k] = v
	}
}

type roundRobinSelector struct {
	servers []string
	i       int
}

func newRoundRobinSelector(servers map[string]string) client.Selector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}
	return &roundRobinSelector{servers: ss}
}

func (s *roundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := s.i
	i = i % len(ss)
	s.i = i + 1

	return ss[i]
}

func (s *roundRobinSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}
