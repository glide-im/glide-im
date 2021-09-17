package route

import (
	"context"
	"errors"
	"github.com/smallnest/rpcx/share"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) *Client {
	return &Client{
		Cli: rpc.NewBaseClient(options),
	}
}

func (c *Client) Unregister(srvId string) error {
	return c.Call(context.Background(), "Unregister", &pb.UnRegisterReq{SrvId: srvId}, &emptypb.Empty{})
}

func (c *Client) Register(param *pb.RegisterRtReq, reply *emptypb.Empty) error {
	return c.Call(context.Background(), "Register", param, reply)
}

func (c *Client) SetTag(srvId, tag, value string) error {
	req := &pb.SetTagReq{
		Tag:   tag,
		SrvId: srvId,
		Value: value,
	}
	return c.Call(context.Background(), "SetTag", req, &emptypb.Empty{})
}

func (c *Client) RemoveTag(srvId, tag string) error {
	return c.Call(context.Background(), "RemoveTag", &pb.ClearTagReq{
		SrvId: srvId,
		Tag:   tag,
	}, &emptypb.Empty{})
}

func (c *Client) Route(ctx context.Context, target string, request, reply interface{}) error {

	split := strings.Split(target, ".")
	if len(split) != 2 {
		return errors.New("参数 target 格式错误, (srvId.func).() 例子: api.Handle")
	}

	var reqParam *anypb.Any
	var err error

	if p, ok := request.(proto.Message); ok {
		reqParam, err = anypb.New(p)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request must be proto.Message")
	}

	routeReq := &pb.RouteReq{
		SrvId:  split[0],
		Fn:     split[1],
		Params: reqParam,
		Extra:  map[string]string{},
	}
	routeReply := &pb.RouteReply{}
	err = c.Call(ctx, "Route", routeReq, routeReply)

	if err != nil {
		return err
	}
	if resp, ok := reply.(proto.Message); ok {
		if !routeReply.GetReply().MessageIs(resp) {
			return errors.New("route reply message not matched to source reply")
		}
		err = routeReply.GetReply().UnmarshalTo(resp)
		if err != nil {
			return err
		}
	}
	return c.Call(ctx, "Route", request, reply)
}

func (c *Client) Route2(target string, request interface{}, reply interface{}) error {
	return c.Route(context.Background(), target, request, reply)
}

func (c *Client) RouteByTag(target, tag string, request, reply interface{}) error {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{ExtraTag: tag})
	return c.Route(ctx, target, request, reply)
}
