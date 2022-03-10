package group_messaging

import (
	"context"
	"go_im/pkg/logger"
)

type groupSelector struct {
	srvs map[string]string
}

func (g *groupSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	for k, v := range g.srvs {
		logger.D("%s: %s", v, k)
		return k
	}
	return ""
}

func (g *groupSelector) UpdateServer(servers map[string]string) {
	g.srvs = servers
}
