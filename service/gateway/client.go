package gateway

import (
	"context"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
	rpc2 "github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
)

type Client struct {
	rpc2.Cli
}

func NewClient(options *rpc2.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = rpc2.NewBaseClient(options)
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

func (c *Client) EnqueueMessage(uid int64, device int64, msg *message.Message) error {

	req := &pb_rpc.EnqueueMessageRequest{
		Uid:     uid,
		Message: msg.GetProtobuf(),
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
	ret := rpc2.NewCtxFrom(context.Background())
	return ret
}

func wrapMessage(msg *message.Message) *message.Message {
	return &message.Message{}
}
