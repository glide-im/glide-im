package route

import (
	"context"
	"errors"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type Client struct {
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	return &Client{
		BaseClient: rpc.NewBaseClient(options),
	}
}

func (c *Client) Invoke2(ctx context.Context, param *pb.RouteReq, reply *pb.Response) error {
	return c.Call("Route", param, reply)
}

func (c *Client) Invoke(target string, request interface{}, reply interface{}) error {
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
	}
	routeReq := &pb.RouteReq{
		SrvId:  split[0],
		Fn:     split[1],
		Params: reqParam,
		Extra:  map[string]string{},
	}
	routeReply := &pb.RouteReply{}
	err = c.Call("Route", routeReq, routeReply)
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
	return nil
}

func (c *Client) Register(param *pb.RegisterRtReq, reply *emptypb.Empty) error {
	return c.Call2(context.Background(), "Register", param, reply)
}
