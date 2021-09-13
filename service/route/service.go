package route

import (
	"context"
	"go_im/pkg/logger"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	*rpc.BaseClient
	name string
}

func newService(options *rpc.ClientOptions) *service {
	return &service{
		BaseClient: rpc.NewBaseClient(options),
		name:       options.Name,
	}
}

func (r *service) route(ctx context.Context, fn string, param interface{}, reply interface{}) error {
	_ = r.Call2(ctx, fn, param, &emptypb.Empty{})
	logger.D("%s.%s", r.name, fn)
	return nil
}
