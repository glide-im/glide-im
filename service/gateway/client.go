package gateway

import (
	"context"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/protobuff/gen/pb_rpc"
	"go_im/service/rpc"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	client.SetInterfaceImpl(ret)
	return ret, nil
}

func (c *Client) ClientSignIn(id int64, uid int64, device int64) error {
	req := &pb_rpc.GatewaySignInRequest{
		Old:    id,
		Uid:    uid,
		Device: device,
	}
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(id, device), "SignIn", req, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) ClientLogout(uid int64, device int64) error {
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(uid, device), "Logout", &pb_rpc.GatewayLogoutRequest{Uid: uid, Device: device}, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) EnqueueMessage(uid int64, device int64, message *message.Message) error {

	req := &pb_rpc.EnqueueMessageRequest{
		Uid: uid,
	}
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(uid, -1), "EnqueueMessage", req, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) isDeviceOnline(uid, device int64) bool {
	return true
}

func (c *Client) allClient() []int64 {
	// TODO
	return nil
}

func getTagContext(uid int64, device int64) context.Context {
	ret := rpc.NewCtxFrom(context.Background())

	return ret
}

func wrapMessage(msg *message.Message) *message.Message {
	return &message.Message{}
}
