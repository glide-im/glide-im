package route

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"go_im/pkg/logger"
	"go_im/service/rpc"
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
		round:    rpc.NewServerSelector(),
		tags:     map[string]string{},
	}
}

func (r *selector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)

	if tag, ok := m[ExtraTag]; ok {
		if path, ok := r.tags[tag]; ok {
			if s, ok := r.services[path]; ok {
				logger.D("route by tag: %s=%s", tag, path)
				return s
			}
		}
	}
	return r.round.Select(ctx, servicePath, serviceMethod, args)
}

func (r *selector) UpdateServer(servers map[string]string) {
	r.round.UpdateServer(servers)
	for k, v := range servers {
		r.services[k] = v
	}
}
