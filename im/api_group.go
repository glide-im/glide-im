package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) CreateGroup(msg *ApiMessage, request *entity.CreateGroupRequest) error {

	group, err := dao.GroupDao.CreateGroup(request.Name, msg.uid)
	if err != nil {
		return err
	}
	chat, err := dao.ChatDao.CreateChat(dao.ChatTypeGroup, group.Gid)
	if err != nil {
		return err
	}

	g := NewGroup(group.Gid, group, chat.Cid, []*dao.GroupMember{})
	GroupManager.AddGroup(g)

	// add self as admin
	owner, err := dao.GroupDao.AddMember(group.Gid, dao.GroupMemberAdmin, msg.uid)
	if err != nil {
		return err
	}
	g.PutMember(owner[0], ClientManager.GetClient(msg.uid).messages)

	// add invited member to group
	if len(request.Member) > 0 {
		nMsg := &ApiMessage{
			uid: msg.uid,
			seq: -1,
		}
		nReq := &entity.AddMemberRequest{
			Gid: group.Gid,
			Uid: request.Member,
		}
		err = m.AddGroupMember(nMsg, nReq)
		if err != nil {
			resp := entity.NewSimpleMessage(-1, entity.RespActionGroupAddMember, "add member failed, "+err.Error())
			ClientManager.EnqueueMessage(msg.uid, resp)
		}
	}

	body := entity.AddedGroupResponse{
		Group:     group,
		GroupChat: chat,
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(body); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupMember(msg *ApiMessage, request *entity.GetGroupMemberRequest) error {

	var ms []*entity.GroupMemberResponse

	g := GroupManager.GetGroup(request.Gid)

	for _, member := range g.members {
		ms = append(ms, &entity.GroupMemberResponse{
			Uid:        member.Uid,
			Nickname:   "",
			RemarkName: member.Remark,
			Type:       member.Type,
			Online:     ClientManager.IsOnline(member.Uid),
		})
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err := resp.SetData(ms); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupInfo(msg *ApiMessage, request *entity.GroupInfoRequest) error {
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "get group info"))
	return nil
}

func (m *groupApi) RemoveMember(msg *ApiMessage, request *entity.RemoveMemberRequest) error {

	for _, uid := range request.Uid {
		GroupManager.GetGroup(request.Gid).Unsubscribe(uid)
		err := dao.GroupDao.RemoveMember(request.Gid, uid)

		if err != nil {
			return err
		}
		notifyResp := entity.NewSimpleMessage(-1, entity.RespActionGroupRemoved, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, notifyResp)
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "remove member success")

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddGroupMember(msg *ApiMessage, request *entity.AddMemberRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	var uids []int64
	for _, uid := range request.Uid {
		// exist member
		if !g.HasMember(uid) {
			uids = append(uids, uid)
		}
	}

	if len(uids) == 0 {
		resp := entity.NewSimpleMessage(msg.seq, entity.RespActionGroupAddMember, "already added")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	members, err := dao.GroupDao.AddMember(request.Gid, dao.GroupMemberUser, uids...)
	if err != nil {
		return err
	}

	for _, member := range members {
		client := ClientManager.GetClient(member.Uid)
		g.PutMember(member, client.messages)

		chat, er := dao.ChatDao.NewUserChat(g.Cid, member.Uid, g.Gid, dao.ChatTypeGroup)
		if er != nil {
			continue
		}
		// notify if member online
		if client != nil {
			client.AddGroup(g)
			chatNotify := entity.NewMessage2(-1, entity.ActionUserNewChat, chat)
			respNotifyMember := entity.NewMessage2(-1, entity.ActionGroupAddMember, g)
			ClientManager.EnqueueMessageMulti(member.Uid, respNotifyMember, chatNotify)
		}
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionGroupAddMember, "add member success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)
	g.Unsubscribe(msg.uid)

	err := dao.GroupDao.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "exit group success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}

func (m *groupApi) JoinGroup(msg *ApiMessage, request *entity.JoinGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	if g == nil {
		resp := entity.NewSimpleMessage(msg.seq, entity.RespActionFailed, "group not exist")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	if g.HasMember(msg.uid) {
		resp := entity.NewSimpleMessage(msg.seq, entity.RespActionFailed, "already joined group")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	_, err := dao.GroupDao.AddMember(request.Gid, 1, msg.uid)
	if err != nil {
		return err
	}

	chat, err := dao.ChatDao.NewUserChat(g.Cid, msg.uid, g.Gid, 2)

	if err != nil {
		_ = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
		return err
	}

	client := ClientManager.GetClient(msg.uid)
	if client != nil {
		client.AddGroup(g)
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(chat); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)

	return nil
}
