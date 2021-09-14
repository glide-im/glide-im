package rpc

import (
	"context"
	"github.com/smallnest/rpcx/share"
)

type HostRouter struct {
	services map[string]string
}

func NewHostRouter() *HostRouter {
	return &HostRouter{services: make(map[string]string)}
}

func (h *HostRouter) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	for k := range h.services {
		return k
	}
	return ""
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

func (h *HostRouter) UpdateServer(servers map[string]string) {
	for k, v := range servers {
		h.services[k] = v
	}
}
