package route

import (
	"context"
	"github.com/smallnest/rpcx/share"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
)

type service struct {
	*rpc.BaseClient
	name string
	tags map[string]string
}

func newService(options *rpc.ClientOptions) *service {
	// unmarshal Any to exactly type
	options.SerializeType = rpc.SerialTypeProtoBuffWrapAny
	options.Selector = newRoundRobinSelector(map[string]string{})

	ret := &service{
		BaseClient: rpc.NewBaseClient(options),
		name:       options.Name,
		tags:       map[string]string{},
	}
	return ret
}

func (r *service) addTag(tag string, value string) {
	r.tags[tag] = value
}

func (r *service) removeTag(tag string) {
	delete(r.tags, tag)
}

func (r *service) route(ctx context.Context, fn string, param *pb.RouteReq, reply *pb.RouteReply) error {
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	tag, ok := m[ExtraTag]
	if ok {
		logger.D("tag: %s", tag)
	}
	_ = r.Call2(ctx, fn, param.GetParams(), reply.GetReply())
	logger.D("%s.%s", r.name, fn)
	return nil
}
