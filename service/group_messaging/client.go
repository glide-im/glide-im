package group_messaging

import (
	"context"
	"go_im/im/group"
	"go_im/im/message"
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
	return ret, nil
}

func (c *Client) DispatchNotifyMessage(gid int64, message *message.GroupNotify) error {
	return c.Call(getContext(gid), "UpdateMember", message, nil)
}

func (c *Client) DispatchMessage(gid int64, action message.Action, message *message.ChatMessage) error {
	return c.Call(getContext(gid), "UpdateMember", message, nil)
}

func (c *Client) UpdateMember(gid int64, update []group.MemberUpdate) error {
	return c.Call(getContext(gid), "UpdateMember", update, nil)
}

func (c *Client) UpdateGroup(gid int64, update group.Update) error {
	return c.Call(getContext(gid), "UpdateMember", update, nil)
}

func getContext(gid int64) context.Context {
	ctx := rpc.NewCtxFrom(context.Background())
	return ctx
}
