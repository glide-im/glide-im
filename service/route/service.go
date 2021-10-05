package route

import (
	"context"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
	"strconv"
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
	switch r.name {
	case "client":
		return r.routeClient(ctx, fn, param, reply)
	case "group":
		return r.routeGroup(ctx, fn, param, reply)
	default:
		return r.Call(ctx, fn, param.GetParams(), reply.GetReply())
	}
}

func (r *service) routeClient(ctx context.Context, fn string, param *pb.RouteReq, reply *pb.RouteReply) error {
	c := rpc.NewCtxFrom(ctx)
	uid, eUid := c.GetReqExtra(ExtraUid)
	device, eDevice := c.GetReqExtra(ExtraDevice)
	if eUid && eDevice {
		uidInt, err := strconv.ParseInt(uid, 10, 64)
		deviceInt, err2 := strconv.ParseInt(device, 10, 64)
		if err == nil && err2 == nil {
			rt := getDeviceRoute(uidInt, deviceInt)
			c.PutReqExtra(ExtraTarget, rt)
		}
	}
	return r.Call(c, fn, param.GetParams(), reply.GetReply())
}

func (r *service) routeGroup(ctx context.Context, fn string, param *pb.RouteReq, reply *pb.RouteReply) error {
	c := rpc.NewCtxFrom(ctx)
	gid, ok := c.GetReqExtra(ExtraGid)
	if ok {
		gidInt, err := strconv.ParseInt(gid, 10, 64)
		if err == nil {
			route := getGroupRoute(gidInt)
			c.PutReqExtra(ExtraTarget, route)
		}
	}
	return r.Call(c, fn, param.GetParams(), reply.GetReply())
}
