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
	// create group chat
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
	// add group to self contacts list
	_, err = dao.UserDao.AddContacts(msg.uid, group.Gid, dao.ContactsTypeGroup, "")
	if err != nil {
		return err
	}
	c := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{},
		Groups: []*entity.GroupResponse{{
			Group:   *g.group,
			Members: g.GetMembers(),
		}},
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

	// create user chat by default
	uc, err := dao.ChatDao.NewUserChat(chat.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, uc))

	// subscribe message
	g.PutMember(owner[0])
	ClientManager.SubscribeGroup(msg.uid, g.Gid)

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

	g := GroupManager.GetGroup(request.Gid)
	ms := make([]*entity.GroupMemberResponse, 0, len(g.members))

	for _, member := range g.members {
		ms = append(ms, &entity.GroupMemberResponse{
			Uid:        member.Uid,
			Nickname:   "",
			RemarkName: member.Remark,
			Type:       member.Type,
			Online:     ClientManager.IsOnline(member.Uid),
		})
	}

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess)
	if err := resp.SetData(ms); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupInfo(msg *ApiMessage, request *entity.GroupInfoRequest) error {

	var groups []*entity.GroupResponse

	for _, gid := range request.Gid {
		g := GroupManager.GetGroup(gid)
		gr := entity.GroupResponse{
			Group:   *g.group,
			Members: g.GetMembers(),
		}
		groups = append(groups, &gr)
	}

	resp := entity.NewMessage2(msg.seq, entity.ActionSuccess, groups)
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) RemoveMember(msg *ApiMessage, request *entity.RemoveMemberRequest) error {

	for _, uid := range request.Uid {
		GroupManager.GetGroup(request.Gid).Unsubscribe(uid)
		err := dao.GroupDao.RemoveMember(request.Gid, uid)

		if err != nil {
			return err
		}
		notifyResp := entity.NewSimpleMessage(-1, entity.ActionGroupRemoveMember, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, notifyResp)
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "remove member success")

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddGroupMember(msg *ApiMessage, request *entity.AddMemberRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	uids := make([]int64, 0, len(request.Uid))
	for _, uid := range request.Uid {
		// member exist
		if !g.HasMember(uid) {
			uids = append(uids, uid)
		}
	}

	if len(uids) == 0 {
		resp := entity.NewSimpleMessage(msg.seq, entity.ActionFailed, "already added")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	// TODO query user info and notify group members, optimize query time
	exist, err2 := dao.UserDao.HasUser(uids...)
	if err2 != nil {
		return err2
	}
	if !exist {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "user does not exist"))
		return nil
	}

	members, err := dao.GroupDao.AddMember(request.Gid, dao.GroupMemberUser, uids...)
	if err != nil {
		return err
	}
	// notify group member update group
	groupNotify := entity.GroupAddMemberResponse{
		Gid:     g.Gid,
		Members: members,
	}
	g.SendMessage(msg.uid, entity.NewMessage2(-1, entity.ActionGroupAddMember, groupNotify))

	for _, member := range members {

		// add group to member's contacts list
		_, e := dao.UserDao.AddContacts(member.Uid, g.Gid, dao.ContactsTypeGroup, "")
		if e != nil {
			continue
		}
		//notify update contacts list
		c := entity.ContactResponse{
			Friends: []*entity.UserInfoResponse{},
			Groups: []*entity.GroupResponse{{
				Group:   *g.group,
				Members: g.GetMembers(),
			}},
		}
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

		// default add user chat
		chat, er := dao.ChatDao.NewUserChat(g.Cid, member.Uid, g.Gid, dao.ChatTypeGroup)
		if er != nil {
			continue
		}
		// member subscribe group message
		g.PutMember(member)
		ClientManager.SubscribeGroup(msg.uid, g.Gid)

		// notify update chat list
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))

	}

	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	GroupManager.UserSignOut(msg.uid, request.Gid)

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

	if g.HasMember(msg.uid) {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "already joined group"))
		return nil
	}

	members, err := dao.GroupDao.AddMember(request.Gid, dao.GroupMemberUser, msg.uid)
	if err != nil {
		return err
	}
	g.PutMember(members[0])

	_, err = dao.UserDao.AddContacts(msg.uid, g.Gid, dao.ContactsTypeGroup, "")
	c := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{},
		Groups: []*entity.GroupResponse{{
			Group:   *g.group,
			Members: g.GetMembers(),
		}},
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

	chat, err := dao.ChatDao.NewUserChat(g.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {
		_ = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))

	ClientManager.SubscribeGroup(msg.uid, g.Gid)
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "join group success"))

	return nil
}
