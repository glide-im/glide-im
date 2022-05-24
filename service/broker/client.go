package broker

import (
	"context"
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error

	options.Selector = newBrokerSelector()
	ret.Cli, err = rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) DispatchNotifyMessage(gid int64, message *message.GroupNotify) error {
	return c.Call("DispatchNotifyMessage", message, &pb_rpc.Response{})
}

func (c *Client) DispatchMessage(gid int64, action message.Action, message *message.ChatMessage) error {
	param := pb_rpc.DispatchGroupChatParam{
		Gid:     gid,
		Action:  string(action),
		Message: message.ChatMessage,
	}
	return c.Call("DispatchMessage", &param, &pb_rpc.Response{})
}

func (c *Client) UpdateMember(gid int64, update []group.MemberUpdate) error {

	var updates []*pb_rpc.MemberUpdateParam
	for _, u := range update {
		up := pb_rpc.MemberUpdateParam{
			Uid:  u.Uid,
			Flag: u.Flag,
		}
		updates = append(updates, &up)
	}
	param := pb_rpc.UpdateMemberParam{
		Gid:     gid,
		Updates: updates,
	}
	return c.Call("UpdateMember", &param, &pb_rpc.Response{})
}

func (c *Client) UpdateGroup(gid int64, update group.Update) error {
	param := pb_rpc.UpdateGroupParam{
		Gid:  gid,
		Flag: update.Flag,
	}
	return c.Call("UpdateGroup", &param, &pb_rpc.Response{})
}

func (c *Client) Call(fn string, param interface{}, replay interface{}) error {
	g, ok := param.(gidParam)
	ctx := context.TODO()
	if ok {
		gid := g.GetGid()
		ctx = context.WithValue(ctx, ctxKeyGid, gid)
	}
	return c.Cli.Call(ctx, fn, param, replay)
}

type gidParam interface {
	GetGid() int64
}
