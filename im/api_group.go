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
	// add group to self contacts list
	c, err := dao.UserDao.AddContacts(msg.uid, group.Gid, dao.ContactsTypeGroup, "")
	if err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, c))

	// create user chat by default
	uc, err := dao.ChatDao.NewUserChat(chat.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {

	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, uc))

	// subscribe message
	client := ClientManager.GetClient(msg.uid)
	g.PutMember(owner[0], client.messages)
	if client != nil {
		client.AddGroup(g)
	}

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

	body := entity.AddedGroupResponse{
		Group:     group,
		GroupChat: chat,
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(msg.seq, entity.ActionSuccess, body))
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

	members, err := dao.GroupDao.AddMember(request.Gid, dao.GroupMemberUser, uids...)
	if err != nil {
		return err
	}

	for _, member := range members {

		// add group to member's contacts list
		_, e := dao.UserDao.AddContacts(member.Uid, g.Gid, dao.ContactsTypeGroup, "")
		if e != nil {
			continue
		}
		//notify update contacts list
		//ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, contacts))

		// default add user chat
		chat, er := dao.ChatDao.NewUserChat(g.Cid, member.Uid, g.Gid, dao.ChatTypeGroup)
		if er != nil {
			continue
		}

		// member subscribe group message
		client := ClientManager.GetClient(member.Uid)
		g.PutMember(member, client.messages)
		client.AddGroup(g)

		// notify update chat list
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))
		if client != nil {
			client.AddGroup(g)
		}
	}

	body := entity.GroupResponse{
		Group:   *g.group,
		Members: g.group.Members,
	}
	r := entity.NewMessage2(-1, entity.ActionGroupJoin, body)
	for _, member := range members {
		ClientManager.EnqueueMessage(member.Uid, r)
	}

	groupNotify := entity.GroupAddMemberResponse{
		Gid:     g.Gid,
		Members: members,
	}
	g.SendMessage(msg.uid, entity.NewMessage2(-1, entity.ActionGroupAddMember, groupNotify))
	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)
	g.Unsubscribe(msg.uid)

	err := dao.GroupDao.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "exit group success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}

func (m *groupApi) JoinGroup(msg *ApiMessage, request *entity.JoinGroupRequest) error {

	client := ClientManager.GetClient(msg.uid)
	g := GroupManager.GetGroup(request.Gid)

	if g == nil {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "group not exist"))
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
	g.PutMember(members[0], client.messages)

	contacts, err := dao.UserDao.AddContacts(msg.uid, g.Gid, dao.ContactsTypeGroup, "")
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, contacts))

	chat, err := dao.ChatDao.NewUserChat(g.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {
		_ = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(-1, entity.ActionUserNewChat, chat))

	client.AddGroup(g)
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.ActionSuccess, "join group success"))

	return nil
}
