package im_service

import (
	"context"
	"errors"
	"go_im/im/client"
	"go_im/im/message"
	rpc2 "go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
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
		return errors.New("im service rpc call error")
	}
	if !resp.Ok {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *Client) ClientLogout(uid int64, device int64) error {
	resp := &pb_rpc.Response{}
	request := &pb_rpc.GatewayLogoutRequest{Uid: uid, Device: device}
	err := c.Call(getTagContext(uid, device), "Logout", request, resp)
	if err != nil {
		return errors.New("im service rpc call error")
	}
	if !resp.Ok {
		return errors.New(resp.Message)
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
		return errors.New("im service rpc call error")
	}
	if !resp.Ok {
		return errors.New(resp.Message)
	}
	return nil
}

func getTagContext(uid int64, device int64) context.Context {
	ret := rpc2.NewCtxFrom(context.Background())
	return ret
}
