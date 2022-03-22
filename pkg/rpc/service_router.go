package rpc

import (
	"context"
	"github.com/smallnest/rpcx/client"
)

// RoundRobinSelector selects servers with roundrobin.
type RoundRobinSelector struct {
	servers []string
	i       int
}

func NewRoundRobinSelector() client.Selector {
	return &RoundRobinSelector{servers: []string{}}
}

func (s *RoundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	return s.SelectNext()
}

func (s *RoundRobinSelector) SelectNext() string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := s.i
	i = i % len(ss)
	s.i = i + 1
	return ss[i]
}

func (s *RoundRobinSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}
