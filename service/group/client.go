package group

import (
	"context"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
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

func NewClientByRouter(srvId string, rtOpts *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = route.NewRouter(srvId, rtOpts)
	if err != nil {
		return nil, err
	}
	group.Manager = ret
	return ret, nil
}

func (c *Client) PutMember(gid int64, mb map[int64]int32) {
	req := &pb.PutMemberRequest{
		Gid:    gid,
		Member: mb,
	}
	resp := &pb.Response{}
	err := c.Call(getContext(gid), "PutMember", req, resp)
	if err != nil {

	}
}

func (c *Client) RemoveMember(gid int64, uid ...int64) error {
	req := &pb.RemoveMemberRequest{
		Gid: gid,
		Uid: uid,
	}
	resp := &pb.Response{}
	err := c.Call(getContext(gid), "RemoveMember", req, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) AddGroup(group *dao.Group, owner int64) {
	req := &pb.AddGroupRequest{
		Group: daoGroup2pbGroup(group),
		Owner: owner,
	}
	resp := &pb.Response{}
	err := c.Call(getContext(group.Gid), "AddGroup", req, resp)
	if err != nil {

	}
}

func (c *Client) UserOnline(uid, gid int64) {
	//resp := &pb.Response{}
	//err := c.Call(context.Background(),"PutMember", req, resp)
	//if err != nil {
	//
	//}
}

func (c *Client) UserOffline(uid, gid int64) {
	//resp := &pb.Response{}
	//err := c.Call(context.Background(),"PutMember", req, resp)
	//if err != nil {
	//
	//}
}

func (c *Client) DispatchNotifyMessage(uid int64, gid int64, message *message.Message) {
	req := &pb.NotifyRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}
	err := c.Call(getContext(gid), "DispatchNotifyMessage", req, resp)
	if err != nil {

	}
}

func (c *Client) DispatchMessage(uid int64, message *message.Message) {
	req := &pb.DispatchMessageRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}
	err := c.Call(context.Background(), "HandleMessage", req, resp)
	if err != nil {
		logger.E("dispatch group message", err)
	}
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
		Data:   msg.Data,
	}
}
