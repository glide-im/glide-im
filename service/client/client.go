package client

import (
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/rpc"
)

type Client struct {
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	ret := &Client{}
	ret.BaseClient = rpc.NewBaseClient(options)
	client.Manager = ret
	return ret
}

// idle function
func (c *Client) ClientConnected(conn conn.Connection) int64 {
	return 0
}

func (c *Client) ClientSignIn(oldUid int64, uid int64, device int64) {
	req := &pb.SignInRequest{
		Old:    oldUid,
		Uid:    uid,
		Device: device,
	}
	resp := &pb.Response{}
	err := c.Call("ClientSignIn", req, resp)
	if err != nil {

	}
}

func (c *Client) UserLogout(uid int64) {
	resp := &pb.Response{}
	err := c.Call("UserLogout", &pb.UidRequest{Uid: uid}, resp)
	if err != nil {

	}
}

func (c *Client) HandleMessage(from int64, message *message.Message) error {
	req := &pb.UidMessageRequest{
		From:    from,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}

	err := c.Call("HandleMessage", req, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) EnqueueMessage(uid int64, message *message.Message) {

	req := &pb.UidMessageRequest{
		From:    uid,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}
	err := c.Call("EnqueueMessage", req, resp)
	if err != nil {

	}
}

func (c *Client) AllClient() []int64 {
	// TODO
	return nil
}

func wrapMessage(msg *message.Message) *pb.Message {
	return &pb.Message{
		Seq:    msg.Seq,
		Action: string(msg.Action),
		Data:   msg.Data,
	}
}
