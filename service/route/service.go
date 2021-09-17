package route

import (
	"context"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
)

type service struct {
	*rpc.BaseClient
	name     string
	selector *selector
}

func newService(options *rpc.ClientOptions) *service {
	// unmarshal Any to exactly type
	options.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	s := newSelector()
	options.Selector = s

	ret := &service{
		BaseClient: rpc.NewBaseClient(options),
		name:       options.Name,
		selector:   s,
	}
	return ret
}

func (r *service) addTag(tag string, value string) {
	r.selector.tags[tag] = value
}

func (r *service) removeTag(tag string) {
	delete(r.selector.tags, tag)
}

func (r *service) route(ctx context.Context, fn string, param *pb.RouteReq, reply *pb.RouteReply) error {
	_ = r.Call(ctx, fn, param.GetParams(), reply.GetReply())
	logger.D("%s.%s", r.name, fn)
	return nil
}
