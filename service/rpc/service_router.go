package rpc

import (
	"context"
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

func (h *HostRouter) UpdateServer(servers map[string]string) {
	for k, v := range servers {
		h.services[k] = v
	}
}
