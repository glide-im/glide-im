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

func newService(options *rpc.ClientOptions) (*service, error) {
	// unmarshal Any to exactly type
	options.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	s := newSelector()
	options.Selector = s

	c, err := rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	ret := &service{
		BaseClient: c,
		name:       options.Name,
		selector:   s,
	}
	return ret, nil
}

func (r *service) addTag(tag string, value string) {
	r.selector.tags[tag] = value
}

func (r *service) removeTag(tag string) {
	delete(r.selector.tags, tag)
}

func (r *service) route(ctx context.Context, fn string, param *pb.RouteReq, reply *pb.RouteReply) error {
	logger.D("%s.%s", r.name, fn)
	return r.Call(ctx, fn, param.GetParams(), reply.GetReply())
}
