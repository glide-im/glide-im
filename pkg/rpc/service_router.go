package rpc

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

type roundRobinSelector struct {
	servers []string
	i       int
}

func NewServerSelector() client.Selector {
	return &roundRobinSelector{servers: []string{}}
}

func (s *roundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	c := NewCtxFrom(ctx)
	c.PutReqExtra("", "")

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

func LogContext(ctx context.Context) {
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	for k, v := range m {
		print("req_mate_data: ", k, ":", v, "\n")
	}
	m2, ok := ctx.Value(share.ResMetaDataKey).(map[string]string)
	if !ok {
		return
	}
	for k, v := range m2 {
		print("res_mate_data: ", k, ":", v, "\n")
	}
}
