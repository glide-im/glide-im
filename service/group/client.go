package group

import (
	"context"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/service/group/pb"
	"go_im/service/rpc"
)

type Client struct {
	rpc pb.GroupServiceClient
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	ret := &Client{}
	ret.BaseClient = rpc.NewBaseClient(options)
	ret.Init(options)
	return ret
}

func (c *Client) PutMember(gid int64, mb *dao.GroupMember) {
	ctx := context.TODO()
	req := &pb.PutMemberRequest{
		Gid:    gid,
		Member: daoMember2pbMember(mb)[0],
	}
	_, err := c.rpc.PutMember(ctx, req)
	if err != nil {

	}
}

func (c *Client) RemoveMember(gid int64, uid ...int64) error {
	ctx := context.TODO()
	_, err := c.rpc.RemoveMember(ctx, &pb.RemoveMemberRequest{
		Gid: gid,
		Uid: uid,
	})
	if err != nil {

	}
	return nil
}

func (c *Client) GetMembers(gid int64) ([]*dao.GroupMember, error) {
	ctx := context.TODO()
	members, err := c.rpc.GetMembers(ctx, &pb.GidRequest{Gid: gid})
	if err != nil {
		return nil, err
	}
	return pbMember2daoMember(members.Members...), err
}

func (c *Client) AddGroup(group *dao.Group, cid int64, owner *dao.GroupMember) {
	ctx := context.TODO()
	_, err := c.rpc.AddGroup(ctx, &pb.AddGroupRequest{
		Group: daoGroup2pbGroup(group),
		Cid:   cid,
		Owner: daoMember2pbMember(owner)[0],
	})
	if err != nil {

	}
}

func (c *Client) GetGroup(gid int64) *dao.Group {
	ctx := context.TODO()
	group, err := c.rpc.GetGroup(ctx, &pb.GidRequest{Gid: gid})
	if err != nil {

	}
	return pbGroup2daoGroup(group)
}

func (c *Client) GetGroupCid(gid int64) int64 {
	ctx := context.TODO()
	cid, err := c.rpc.GetGroupCid(ctx, &pb.GidRequest{Gid: gid})
	if err != nil {
		return 0
	}
	return cid.GetCid()
}

func (c *Client) HasMember(gid int64, uid int64) bool {
	ctx := context.TODO()
	member, err := c.rpc.HasMember(ctx, &pb.HasMemberRequest{
		Gid: gid,
		Uid: uid,
	})
	if err != nil {
		return false
	}
	return member.GetHas()
}

func (c *Client) DispatchNotifyMessage(uid int64, gid int64, message *message.Message) {
	ctx := context.TODO()
	_, err := c.rpc.DispatchNotifyMessage(ctx, &pb.NotifyRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	})
	if err != nil {
		return
	}
}

func (c *Client) DispatchMessage(uid int64, message *message.Message) error {
	ctx := context.TODO()
	_, err := c.rpc.DispatchMessage(ctx, &pb.DispatchMessageRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	})
	return err
}

func wrapMessage(msg *message.Message) *pb.Message {
	return &pb.Message{
		Seq:    msg.Seq,
		Action: string(msg.Action),
		Data:   msg.Data,
	}
}

func (c *Client) Run() error {
	err := c.Connect()
	c.rpc = pb.NewGroupServiceClient(c.Conn)
	return err
}
