package gateway

import (
	"context"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/protobuff/pb_rpc"
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
	client.Manager = ret
	return ret, nil
}

// ClientConnected idle function
func (c *Client) ClientConnected(conn conn.Connection) int64 {
	return 0
}

func (c *Client) AddClient(uid int64, cs client.IClient) {}

func (c *Client) ClientSignIn(id int64, uid int64, device int64) {
	req := &pb_rpc.GatewaySignInRequest{
		Old:    id,
		Uid:    uid,
		Device: device,
	}
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(id, device), "SignIn", req, resp)
	if err != nil {

	}
}

func (c *Client) ClientLogout(uid int64, device int64) {
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(uid, device), "Logout", &pb_rpc.GatewayLogoutRequest{Uid: uid, Device: device}, resp)
	if err != nil {

	}
}

func (c *Client) EnqueueMessage(uid int64, device int64, message *message.Message) {

	req := &pb_rpc.EnqueueMessageRequest{
		Uid: uid,
	}
	resp := &pb_rpc.Response{}
	err := c.Call(getTagContext(uid, -1), "EnqueueMessage", req, resp)
	if err != nil {

	}
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
