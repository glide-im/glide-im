package group_messaging

import (
	"context"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
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
	return c.Call(getContext(gid), "DispatchNotifyMessage", message, &pb_rpc.Response{})
}

func (c *Client) DispatchMessage(gid int64, action message.Action, message *message.ChatMessage) error {
	param := pb_rpc.DispatchGroupChatParam{
		Gid:     gid,
		Action:  string(action),
		Message: message.ChatMessage,
	}
	return c.Call(getContext(gid), "DispatchMessage", &param, &pb_rpc.Response{})
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
	return c.Call(getContext(gid), "UpdateMember", &param, &pb_rpc.Response{})
}

func (c *Client) UpdateGroup(gid int64, update group.Update) error {
	param := pb_rpc.UpdateGroupParam{
		Gid:  gid,
		Flag: update.Flag,
	}
	return c.Call(getContext(gid), "UpdateGroup", &param, &pb_rpc.Response{})
}

func getContext(gid int64) context.Context {
	ctx := rpc.NewCtxFrom(context.Background())
	return ctx
}
