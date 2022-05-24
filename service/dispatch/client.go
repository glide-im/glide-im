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

func (c *Client) DispatchGatewayDirect(uid int64, gateway string, m *pb_rpc.NSQGatewayMessage) error {
	any, err := anypb.New(m)
	if err != nil {
		return err
	}
	request := &pb_rpc.DispatchRequest{
		SrvName:  "gateway",
		Id:       uid,
		Data:     any,
		Direct:   true,
		RouteVal: gateway,
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

func (c *Client) GetUserGateway(uid int64) error {
	request := &pb_rpc.UidRequest{
		Uid: uid,
	}
	ctx := contextOfUidHashRoute(uid)
	return c.Call(ctx, "GetUserGateway", request, &pb_rpc.Response{})
}

func (c *Client) ClientSignIn(oldUid int64, uid int64, device int64) error {
	data, _ := anypb.New(&pb_rpc.GatewaySignInRequest{
		Old:    oldUid,
		Uid:    uid,
		Device: device,
	})
	nsqMsg := &pb_rpc.NSQGatewayMessage{
		Operate: pb_rpc.NSQGatewayMessage_LOGIN,
		Params:  data,
	}

	// TODO
	//err := c.UpdateGatewayRoute(uid, "1")
	//if err != nil {
	//	return err
	//}

	return c.DispatchGateway(uid, nsqMsg)
}

func (c *Client) ClientLogout(uid int64, device int64) error {
	data, _ := anypb.New(&pb_rpc.GatewayLogoutRequest{
		Uid:    uid,
		Device: device,
	})
	nsqMsg := &pb_rpc.NSQGatewayMessage{
		Operate: pb_rpc.NSQGatewayMessage_LOGOUT,
		Params:  data,
	}
	return c.DispatchGateway(uid, nsqMsg)
}

func (c *Client) EnqueueMessage(uid int64, device int64, msg *message.Message) error {
	data, _ := anypb.New(&pb_rpc.EnqueueMessageRequest{
		Uid:     uid,
		Message: msg.GetProtobuf(),
	})
	nsqMsg := &pb_rpc.NSQGatewayMessage{
		Operate: pb_rpc.NSQGatewayMessage_PUSH_MSG,
		Params:  data,
	}
	return c.DispatchGateway(uid, nsqMsg)
}
