package dispatch

import (
	"go_im/im/message"
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
	options.Selector = newSelector()
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
	ctx := contextOfUidHashRoute(uid)
	return c.Call(ctx, "Dispatch", request, &pb_rpc.Response{})
}

func (c *Client) UpdateGatewayRoute(uid int64, node string) error {
	request := &pb_rpc.UpdateRouteRequest{
		SrvName: "gateway",
		Id:      uid,
		Node:    node,
	}
	ctx := contextOfUidHashRoute(uid)
	return c.Call(ctx, "UpdateRoute", request, &pb_rpc.Response{})
}

func (c *Client) ClientSignIn(oldUid int64, uid int64, device int64) error {

	return nil
}

func (c *Client) ClientLogout(uid int64, device int64) error {

	return nil
}

func (c *Client) EnqueueMessage(uid int64, device int64, message *message.Message) error {

	return nil
}
