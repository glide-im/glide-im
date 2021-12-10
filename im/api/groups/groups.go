package groups

import (
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/groupdao"
	"go_im/im/group"
	"go_im/im/message"
)

type Interface interface {
}

type GroupApi struct {
}

func (m *GroupApi) CreateGroup(ctx *route.Context, request *CreateGroupRequest) error {

	dbGroup, err := groupdao.Dao.CreateGroup(request.Name, 1)
	if err != nil {
		return comm.NewDbErr(err)
	}
	err = groupdao.Dao.AddMember(dbGroup.Gid, ctx.Uid, MemberTypeAdmin, MemberFlagDefault)
	if err != nil {
		return comm.NewDbErr(err)
	}
	err = groupdao.Dao.AddMembers(dbGroup.Gid, MemberFlagDefault, MemberTypeNormal, request.Member...)
	if err != nil {
		return comm.NewDbErr(err)
	}
	err = apidep.GroupManager.CreateGroup(dbGroup.Gid)
	if err != nil {
		return comm.NewUnexpectedErr("create group failed", err)
	}
	err = apidep.GroupManager.PutMember(dbGroup.Gid, request.Member)
	if err != nil {
		return comm.NewUnexpectedErr("add group member failed", err)
	}
	err = apidep.GroupManager.UpdateMember(dbGroup.Gid, ctx.Uid, group.FlagMemberSetAdmin)
	if err != nil {
		return comm.NewUnexpectedErr("create group failed", err)
	}
	n := message.NewMessage(0, comm.ActionInviteToGroup, InviteGroupMessage{Gid: dbGroup.Gid})
	for _, uid := range request.Member {
		apidep.SendMessageIfOnline(uid, 0, n)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, CreateGroupResponse{Gid: dbGroup.Gid}))
	return nil
}

func (m *GroupApi) GetGroupMember(ctx *route.Context, request *GetGroupMemberRequest) error {

	mbs, err := groupdao.Dao.GetMembers(request.Gid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	ms := make([]*GroupMemberResponse, 0, len(mbs))
	for _, member := range mbs {
		ms = append(ms, &GroupMemberResponse{
			Uid:        member.Uid,
			RemarkName: member.Remark,
			Type:       member.Type,
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ms))
	return nil
}

func (m *GroupApi) GetGroupInfo(ctx *route.Context, request *GroupInfoRequest) error {
	// TODO 2021-12-9 21:55:07
	return nil
}

func (m *GroupApi) RemoveMember(ctx *route.Context, request *RemoveMemberRequest) error {
	// TODO 2021-12-9 21:55:01
	return nil
}

func (m *GroupApi) AddGroupMember(ctx *route.Context, request *AddMemberRequest) error {
	err := addGroupMemberDb(request.Gid, ctx.Uid, MemberFlagDefault)
	if err != nil {
		return err
	}
	err = apidep.GroupManager.PutMember(request.Gid, []int64{ctx.Uid})
	if err != nil {
		return comm.NewUnexpectedErr("add group failed", err)
	}
	for _, i := range request.Uid {
		n := message.NewMessage(0, comm.ActionInviteToGroup, InviteGroupMessage{Gid: request.Gid})
		apidep.SendMessageIfOnline(i, 0, n)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return nil
}

func (m *GroupApi) ExitGroup(ctx *route.Context, request *ExitGroupRequest) error {

	err := apidep.GroupManager.RemoveMember(request.Gid, ctx.Uid)
	if err != nil {
		return comm.NewUnexpectedErr("exit group failed", err)
	}
	err = groupdao.Dao.RemoveMember(request.Gid, ctx.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	resp := message.NewMessage(ctx.Seq, comm.ActionSuccess, " group success")
	ctx.Response(resp)
	return err
}

func (m *GroupApi) JoinGroup(ctx *route.Context, request *JoinGroupRequest) error {
	err := addGroupMemberDb(request.Gid, ctx.Uid, MemberFlagDefault)
	if err != nil {
		return err
	}
	err = apidep.GroupManager.PutMember(request.Gid, []int64{ctx.Uid})
	if err != nil {
		return comm.NewUnexpectedErr("add group failed", err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return nil
}

func addGroupMemberDb(gid int64, uid int64, typ int64) error {
	hasGroup, err := groupdao.Dao.HasGroup(gid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if !hasGroup {
		return ErrGroupNotExit
	}
	hasMember, err := groupdao.Dao.HasMember(gid, uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if hasMember {
		return ErrMemberAlreadyExist
	}
	err = groupdao.Dao.AddMember(gid, uid, typ, MemberFlagDefault)
	if err != nil {
		return comm.NewDbErr(err)
	}
	return nil
}
