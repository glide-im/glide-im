package route

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"go_im/pkg/logger"
)

type selector struct {
	services map[string]string
	round    client.Selector
	tags     map[string]string
}

func newSelector() *selector {
	s := map[string]string{}
	return &selector{
		services: s,
		round:    newRoundRobinSelector(),
		tags:     map[string]string{},
	}
}

func (r *selector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	tag, ok := m[ExtraTag]
	if ok {
		v, ok := r.tags[tag]
		if ok {
			s, ok := r.services[v]
			if ok {
				logger.D("route by tag: %s=%s", tag, v)
				return s
			}
		}
	}
	for k := range r.services {
		return k
	}
	return ""
}

func (r *selector) UpdateServer(servers map[string]string) {
	r.round.UpdateServer(servers)
	for k, v := range servers {
		r.services[k] = v
	}
}

type roundRobinSelector struct {
	servers []string
	i       int
}

func newRoundRobinSelector() client.Selector {
	return &roundRobinSelector{servers: []string{}}
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
