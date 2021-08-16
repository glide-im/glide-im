package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) CreateGroup(msg *ApiMessage, request *entity.CreateGroupRequest) error {

	group, err := GroupManager.CreateGroup(request.Name, msg.uid)
	if err != nil {
		return err
	}
	members, _ := GroupManager.GetMembers(group.Gid)
	c := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{},
		Groups: []*entity.GroupResponse{{
			Group:   *group.group,
			Members: members,
		}},
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

	// create user chat by default
	uc, err := dao.ChatDao.NewUserChat(group.Cid, msg.uid, group.Gid, dao.ChatTypeGroup)
	if err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, uc))
	ClientManager.AddGroup(msg.uid, group.Gid)

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
			resp := entity.NewSimpleMessage(-1, entity.ActionFailed, "add member failed, "+err.Error())
			ClientManager.EnqueueMessage(msg.uid, resp)
		}
	}

	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(msg.seq, entity.ActionSuccess, "create group success"))
	return nil
}

func (m *groupApi) GetGroupMember(msg *ApiMessage, request *entity.GetGroupMemberRequest) error {

	members, err := GroupManager.GetMembers(request.Gid)
	if err != nil {
		return err
	}

	ms := make([]*entity.GroupMemberResponse, 0, len(members))
	for _, member := range members {
		ms = append(ms, &entity.GroupMemberResponse{
			Uid:        member.Uid,
			Nickname:   "",
			RemarkName: member.Remark,
			Type:       member.Type,
			Online:     ClientManager.IsOnline(member.Uid),
		})
	}

	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(msg.seq, entity.ActionSuccess, ms))
	return nil
}

func (m *groupApi) GetGroupInfo(msg *ApiMessage, request *entity.GroupInfoRequest) error {

	var groups []*entity.GroupResponse

	for _, gid := range request.Gid {
		g := GroupManager.GetGroup(gid)
		members, err := GroupManager.GetMembers(gid)
		if err != nil {
			return err
		}
		gr := entity.GroupResponse{
			Group:   *g.group,
			Members: members,
		}
		groups = append(groups, &gr)
	}

	resp := entity.NewMessage2(msg.seq, entity.ActionSuccess, groups)
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) RemoveMember(msg *ApiMessage, request *entity.RemoveMemberRequest) error {

	for _, uid := range request.Uid {
		err := dao.GroupDao.RemoveMember(request.Gid, uid)
		if err != nil {
			return err
		}
		ClientManager.RemoveGroup(uid, request.Gid)
		notifyResp := entity.NewSimpleMessage(-1, entity.ActionGroupRemoveMember, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, notifyResp)
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "remove member success")

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddGroupMember(msg *ApiMessage, request *entity.AddMemberRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	members, err := GroupManager.AddGroupMember(g.Gid, request.Uid...)
	if err != nil {
		return err
	}

	// notify group member update group
	groupNotify := entity.GroupAddMemberResponse{
		Gid:     g.Gid,
		Members: members,
	}
	GroupManager.DispatchNotifyMessage(msg.uid, g.Gid, entity.NewMessage2(-1, entity.ActionGroupAddMember, groupNotify))

	for _, member := range members {

		// add group to member's contacts list
		_, e := dao.UserDao.AddContacts(member.Uid, g.Gid, dao.ContactsTypeGroup, "")
		if e != nil {
			_ = GroupManager.RemoveMember(request.Gid, member.Uid)
			continue
		}
		ms, err := GroupManager.GetMembers(request.Gid)
		if err != nil {
			return err
		}
		//notify update contacts list
		c := entity.ContactResponse{
			Friends: []*entity.UserInfoResponse{},
			Groups: []*entity.GroupResponse{{
				Group:   *g.group,
				Members: ms,
			}},
		}
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

		// default add user chat
		chat, er := dao.ChatDao.NewUserChat(g.Cid, member.Uid, g.Gid, dao.ChatTypeGroup)
		if er != nil {
			continue
		}
		// member subscribe group message
		GroupManager.SubscribeGroup(g.Gid, member)
		ClientManager.AddGroup(msg.uid, g.Gid)

		// notify update chat list
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))
	}

	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	ClientManager.RemoveGroup(msg.uid, request.Gid)

	err := dao.GroupDao.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "exit group success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}

func (m *groupApi) JoinGroup(msg *ApiMessage, request *entity.JoinGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	if g == nil {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "group does not exist"))
		return nil
	}

	_, err := GroupManager.AddGroupMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}

	_, err = dao.UserDao.AddContacts(msg.uid, g.Gid, dao.ContactsTypeGroup, "")

	members, err := GroupManager.GetMembers(request.Gid)
	if err != nil {
		return err
	}

	c := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{},
		Groups: []*entity.GroupResponse{{
			Group:   *g.group,
			Members: members,
		}},
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

	chat, err := dao.ChatDao.NewUserChat(g.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {
		_ = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))

	ClientManager.AddGroup(msg.uid, g.Gid)
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "join group success"))

	return nil
}
