package route

import (
	"context"
	"go_im/pkg/logger"
	"go_im/service/rpc"
)

type service struct {
	*rpc.BaseClient
	name string
}

func newService(options *rpc.ClientOptions) *service {
	// unmarshal Any to exactly type
	options.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	return &service{
		BaseClient: rpc.NewBaseClient(options),
		name:       options.Name,
	}
}

func (r *service) route(ctx context.Context, fn string, param interface{}, reply interface{}) error {
	_ = r.Call2(ctx, fn, param, reply)
	logger.D("%s.%s", r.name, fn)
	return nil
}
