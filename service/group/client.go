package group

import (
	"go_im/im/dao"
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
	return ret
}

func (c *Client) PutMember(gid int64, mb *dao.GroupMember) {
	req := &pb.PutMemberRequest{
		Gid:    gid,
		Member: daoMember2pbMember(mb)[0],
	}
	resp := &pb.Response{}
	err := c.Call("PutMember", req, resp)
	if err != nil {

	}
}

func (c *Client) RemoveMember(gid int64, uid ...int64) error {
	req := &pb.RemoveMemberRequest{
		Gid: gid,
		Uid: uid,
	}
	resp := &pb.Response{}
	err := c.Call("RemoveMember", req, resp)
	if err != nil {

	}
	return nil
}

func (c *Client) GetMembers(gid int64) ([]*dao.GroupMember, error) {
	req := &pb.GidRequest{Gid: gid}
	resp := &pb.GetMembersResponse{}
	err := c.Call("GetMembers", req, resp)
	if err != nil {

	}
	return pbMember2daoMember(resp.Members...), err
}

func (c *Client) AddGroup(group *dao.Group, cid int64, owner *dao.GroupMember) {
	req := &pb.AddGroupRequest{
		Group: daoGroup2pbGroup(group),
		Cid:   cid,
		Owner: daoMember2pbMember(owner)[0],
	}
	resp := &pb.Response{}
	err := c.Call("AddGroup", req, resp)
	if err != nil {

	}
}

func (c *Client) GetGroup(gid int64) *dao.Group {
	req := &pb.GidRequest{Gid: gid}
	resp := &pb.Group{}
	err := c.Call("GetGroup", req, resp)
	if err != nil {

	}
	return pbGroup2daoGroup(resp)
}

func (c *Client) UserOnline(uid, gid int64) {
	//resp := &pb.Response{}
	//err := c.Call("PutMember", req, resp)
	//if err != nil {
	//
	//}
}

func (c *Client) UserOffline(uid, gid int64) {
	//resp := &pb.Response{}
	//err := c.Call("PutMember", req, resp)
	//if err != nil {
	//
	//}
}

func (c *Client) GetGroupCid(gid int64) int64 {
	req := &pb.GidRequest{Gid: gid}
	resp := &pb.GetCidResponse{}
	err := c.Call("GetGroupCid", req, resp)
	if err != nil {

	}
	return resp.GetCid()
}

func (c *Client) HasMember(gid int64, uid int64) bool {
	req := &pb.HasMemberRequest{
		Gid: gid,
		Uid: uid,
	}
	resp := &pb.HasMemberResponse{}
	err := c.Call("HasMember", req, resp)
	if err != nil {

	}
	return resp.GetHas()
}

func (c *Client) DispatchNotifyMessage(uid int64, gid int64, message *message.Message) {
	req := &pb.NotifyRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}
	err := c.Call("DispatchNotifyMessage", req, resp)
	if err != nil {

	}
}

func (c *Client) DispatchMessage(uid int64, message *message.Message) error {
	req := &pb.DispatchMessageRequest{
		Uid:     uid,
		Message: wrapMessage(message),
	}
	resp := &pb.Response{}
	err := c.Call("HandleMessage", req, resp)
	if err != nil {

	}
	return err
}

func wrapMessage(msg *message.Message) *pb.Message {
	return &pb.Message{
		Seq:    msg.Seq,
		Action: string(msg.Action),
		Data:   msg.Data,
	}
}
