package groups

import (
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/common"
	"go_im/im/dao/groupdao"
	"go_im/im/dao/msgdao"
	"go_im/im/dao/userdao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
	"time"
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
	_, err = msgdao.GroupMsgDaoImpl.CreateGroupMessageState(dbGroup.Gid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	err = userdao.Dao.AddContacts(ctx.Uid, dbGroup.Gid, userdao.ContactsTypeGroup)
	if err != nil {
		return comm.NewDbErr(err)
	}

	err = groupdao.Dao.AddMember(dbGroup.Gid, ctx.Uid, groupdao.GroupMemberTypeOwner, groupdao.GroupFlagDefault)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//err = groupdao.Dao.AddMembers(dbGroup.Gid, MemberFlagDefault, MemberTypeNormal, request.Member...)
	//if err != nil {
	//	return comm.NewDbErr(err)
	//}
	err = apidep.GroupManager.CreateGroup(dbGroup.Gid)
	if err != nil {
		return comm.NewUnexpectedErr("create group failed", err)
	}
	err = apidep.GroupManager.PutMember(dbGroup.Gid, []int64{ctx.Uid})
	if err != nil {
		return comm.NewUnexpectedErr("add group member failed", err)
	}
	err = apidep.GroupManager.UpdateMember(dbGroup.Gid, ctx.Uid, group.FlagMemberSetAdmin)
	if err != nil {
		return comm.NewUnexpectedErr("create group failed", err)
	}
	//n := message.NewMessage(0, comm.ActionInviteToGroup, InviteGroupMessage{Gid: dbGroup.Gid})
	//for _, uid := range request.Member {
	//	apidep.SendMessageIfOnline(uid, 0, n)
	//}
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
			Type:       int(member.Type),
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ms))
	return nil
}

func (m *GroupApi) GetGroupInfo(ctx *route.Context, request *GroupInfoRequest) error {
	groups, err := groupdao.Dao.GetGroups(request.Gid...)
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, groups))
	return nil
}

func (m *GroupApi) RemoveMember(ctx *route.Context, request *RemoveMemberRequest) error {
	typ, err := groupdao.Dao.GetMemberType(request.Gid, ctx.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if typ == groupdao.GroupMemberTypeAdmin || typ == groupdao.GroupMemberTypeOwner {
		//goland:noinspection GoPreferNilSlice
		notFind := []int64{}
		for _, id := range request.Uid {
			err = userdao.ContactsDao.DelContacts(id, request.Gid, userdao.ContactsTypeGroup)
			if err != nil {
				return comm.NewDbErr(err)
			}
			err = groupdao.Dao.RemoveMember(request.Gid, id)
			if err == common.ErrNoRecordFound {
				notFind = append(notFind, id)
				continue
			} else if err != nil {
				return comm.NewDbErr(err)
			}
			err = dispatchGroupNotify(request.Gid, message.GroupNotifyTypeMemberRemoved, id)
			if err != nil {
				logger.E("remove member error:%v", err)
			}
		}
		if len(notFind) == 0 {
			ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
		} else {
			ctx.Response(message.NewMessage(ctx.Seq, comm.ActionFailed, notFind))
		}
	} else {
		return ErrGroupNotExit
	}
	return nil
}

func (m *GroupApi) AddGroupMember(ctx *route.Context, request *AddMemberRequest) error {
	for _, uid := range request.Uid {
		err := addGroupMemberDb(request.Gid, uid, groupdao.GroupMemberNormal)
		if err != nil {
			return err
		}
		err = userdao.ContactsDao.AddContacts(uid, request.Gid, userdao.ContactsTypeGroup)
		if err != nil {
			return comm.NewDbErr(err)
		}
		n := message.NewMessage(0, message.ActionNotifyNewContact, comm.NewContactMessage{
			FromId:   ctx.Uid,
			FromType: 0,
			Id:       request.Gid,
			Type:     userdao.ContactsTypeGroup,
		})
		apidep.SendMessageIfOnline(uid, 0, n)
		err = apidep.GroupManager.PutMember(request.Gid, []int64{uid})
		if err != nil {
			return comm.NewUnexpectedErr("add group failed", err)
		}
		err = dispatchGroupNotify(request.Gid, message.GroupNotifyTypeMemberAdded, uid)
		if err != nil {
			logger.E("notify add group member error: %v", err)
		}
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
	err = dispatchGroupNotify(request.Gid, message.GroupNotifyTypeMemberRemoved, ctx.Uid)
	if err != nil {
		logger.E("exit group error: %v", err)
	}
	resp := message.NewMessage(ctx.Seq, comm.ActionSuccess, " group success")
	ctx.Response(resp)
	return err
}

func (m *GroupApi) JoinGroup(ctx *route.Context, request *JoinGroupRequest) error {

	isC, err := userdao.ContactsDao.HasContacts(ctx.Uid, request.Gid, userdao.ContactsTypeGroup)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if isC {
		return ErrMemberAlreadyExist
	}
	// TODO 2021-11-29 use transaction
	err = userdao.ContactsDao.AddContacts(ctx.Uid, request.Gid, userdao.ContactsTypeGroup)
	if err != nil {
		return comm.NewDbErr(err)
	}

	err = addGroupMemberDb(request.Gid, ctx.Uid, groupdao.GroupMemberNormal)
	if err != nil {
		return err
	}
	err = apidep.GroupManager.PutMember(request.Gid, []int64{ctx.Uid})
	if err != nil {
		return comm.NewUnexpectedErr("add group failed", err)
	}
	err = dispatchGroupNotify(request.Gid, message.GroupNotifyTypeMemberAdded, ctx.Uid)
	if err != nil {
		logger.E("join group error:%v", err)
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
	err = groupdao.Dao.AddMember(gid, uid, typ, groupdao.GroupFlagDefault)
	if err != nil {
		return comm.NewDbErr(err)
	}
	return nil
}

func dispatchGroupNotify(gid int64, typ int64, uid int64) error {
	id, err := msgdao.GetMessageID()
	if err != nil {
		logger.E("get message id error:%v", err)
		return err
	}
	notify := message.GroupNotify{
		Mid:       id,
		Gid:       gid,
		Timestamp: time.Now().Unix(),
		Type:      typ,
		Data:      &message.GroupNotifyMemberAdded{Uid: []int64{uid}},
	}
	return apidep.GroupManager.DispatchNotifyMessage(gid, &notify)
}
