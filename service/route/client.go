package route

import (
	"context"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	return &Client{
		BaseClient: rpc.NewBaseClient(options),
	}
}

func (c *Client) Route(ctx context.Context, param *pb.RouteReq, reply *pb.Response) error {
	return c.Call("Route", param, reply)
}

func (c *Client) Register(param *pb.RegisterRtReq, reply *emptypb.Empty) error {
	return c.Call2(context.Background(), "Register", param, reply)
}
