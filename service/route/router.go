package route

import (
	"context"
	"github.com/pkg/errors"
	"go_im/service/rpc"
)

type RouterCli struct {
	rt    *Client
	srvId string
}

func NewRouter(srvId string, routeOpts *rpc.ClientOptions) *RouterCli {
	return &RouterCli{
		rt:    NewClient(routeOpts),
		srvId: srvId,
	}
}

func (r *RouterCli) Call(ctx context.Context, fn string, request, reply interface{}) error {
	path := r.srvId + "." + fn
	return r.rt.Route(ctx, path, request, reply)
}

func (r *RouterCli) Broadcast(fn string, request, reply interface{}) error {
	return errors.New("broadcast on proxy mode is unsupported")
}

func (r *RouterCli) Run() error {
	return r.rt.Run()
}

func (r *RouterCli) Close() error {
	return r.rt.Close()
}
