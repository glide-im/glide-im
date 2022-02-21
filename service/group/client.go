package group

import (
	"context"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/route"
	"go_im/service/rpc"
	"strconv"
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
	return ret, nil
}

func (c *Client) UpdateMember(gid int64, update []group.MemberUpdate) error {
	return c.Call(getContext(gid), "UpdateMember", nil, nil)
}

func (c *Client) UpdateGroup(gid int64, update group.Update) error {
	return c.Call(getContext(gid), "UpdateMember", nil, nil)
}

func (c *Client) DispatchNotifyMessage(gid int64, message *message.Message) error {
	return c.Call(getContext(gid), "UpdateMember", nil, nil)
}

func (c *Client) DispatchMessage(gid int64, message *message.ChatMessage) error {
	return c.Call(getContext(gid), "UpdateMember", nil, nil)
}

func NewClientByRouter(srvId string, rtOpts *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = route.NewRouter(srvId, rtOpts)
	if err != nil {
		return nil, err
	}
	//group.Manager = ret
	return ret, nil
}

func getContext(gid int64) context.Context {
	ctx := rpc.NewCtxFrom(context.Background())
	ctx.PutReqExtra(route.ExtraGid, strconv.FormatInt(gid, 10))
	return ctx
}

func wrapMessage(msg *message.Message) *pb.Message {
	return &pb.Message{
		Seq:    msg.Seq,
		Action: string(msg.Action),
		Data:   "",
	}
}
