package dispatch

import (
	"context"
	rpc2 "go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type Client struct {
	rpc2.Cli
}

func NewClient(options *rpc2.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	options.Selector = &dispatchSelector{}
	ret.Cli, err = rpc2.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) DispatchGateway(uid int64, m *pb_rpc.NSQGatewayMessage) error {
	any, err := anypb.New(m)
	if err != nil {
		return err
	}
	request := &pb_rpc.DispatchRequest{
		SrvName: "gateway",
		Id:      uid,
		Data:    any,
	}
	return c.Call(context.Background(), "Dispatch", request, &pb_rpc.Response{})
}

func (c *Client) UpdateGatewayRoute(uid int64, node string) error {
	request := &pb_rpc.UpdateRouteRequest{
		SrvName: "gateway",
		Id:      uid,
		Node:    node,
	}
	return c.Call(context.Background(), "UpdateRoute", request, &pb_rpc.Response{})
}
