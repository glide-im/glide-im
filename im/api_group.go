package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) CreateGroup(msg *ApiMessage, request *entity.CreateGroupRequest) error {

	group, err := m.createGroup(request.Name, msg.uid)
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
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(-1, entity.ActionUserAddFriend, c))

	// create user chat by default
	uc, err := dao.ChatDao.NewUserChat(group.Cid, msg.uid, group.Gid, dao.ChatTypeGroup)
	if err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(-1, entity.ActionUserNewChat, uc))

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
			resp := entity.NewMessage(-1, entity.ActionFailed, "add member failed, "+err.Error())
			ClientManager.EnqueueMessage(msg.uid, resp)
		}
	}

	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(msg.seq, entity.ActionSuccess, "create group success"))
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

	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(msg.seq, entity.ActionSuccess, ms))
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

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess, groups)
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) RemoveMember(msg *ApiMessage, request *entity.RemoveMemberRequest) error {

	for _, uid := range request.Uid {
		err := dao.GroupDao.RemoveMember(request.Gid, uid)
		if err != nil {
			return err
		}
		_ = GroupManager.RemoveMember(request.Gid, uid)
		notifyResp := entity.NewMessage(-1, entity.ActionGroupRemoveMember, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, notifyResp)
	}

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess, "remove member success")

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddGroupMember(msg *ApiMessage, request *entity.AddMemberRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	members, err := m.addGroupMember(g.Gid, request.Uid...)
	if err != nil {
		return err
	}

	// notify group member update group
	groupNotify := entity.GroupAddMemberResponse{
		Gid:     g.Gid,
		Members: members,
	}
	GroupManager.DispatchNotifyMessage(msg.uid, g.Gid, entity.NewMessage(-1, entity.ActionGroupAddMember, groupNotify))

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
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage(-1, entity.ActionUserAddFriend, c))

		// default add user chat
		chat, er := dao.ChatDao.NewUserChat(g.Cid, member.Uid, g.Gid, dao.ChatTypeGroup)
		if er != nil {
			continue
		}
		// member subscribe group message
		GroupManager.PutMember(g.Gid, member)

		// notify update chat list
		ClientManager.EnqueueMessage(member.Uid, entity.NewMessage(-1, entity.ActionUserNewChat, chat))
	}

	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	err := GroupManager.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}

	err = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewMessage(msg.seq, entity.ActionSuccess, "exit group success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}

func (m *groupApi) JoinGroup(msg *ApiMessage, request *entity.JoinGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	if g == nil {
		return errors.New("group does not exist")
	}

	ms, err := m.addGroupMember(request.Gid, msg.uid)
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
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(-1, entity.ActionUserAddFriend, c))

	chat, err := dao.ChatDao.NewUserChat(g.Cid, msg.uid, g.Gid, dao.ChatTypeGroup)
	if err != nil {
		_ = dao.GroupDao.RemoveMember(request.Gid, msg.uid)
		return err
	}
	GroupManager.PutMember(g.Gid, ms[0])
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(-1, entity.ActionUserNewChat, chat))
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage(msg.seq, entity.ActionSuccess, "join group success"))
	return nil
}

func (m *groupApi) createGroup(name string, uid int64) (*Group, error) {

	group, err := dao.GroupDao.CreateGroup(name, uid)
	if err != nil {
		return nil, err
	}
	// create group chat
	chat, err := dao.ChatDao.CreateChat(dao.ChatTypeGroup, group.Gid)
	if err != nil {
		// TODO undo
		return nil, err
	}
	g := NewGroup(group.Gid, group, chat.Cid, []*dao.GroupMember{})

	owner, err := dao.GroupDao.AddMember(group.Gid, dao.GroupMemberAdmin, uid)
	if err != nil {
		// TODO undo create group
		return nil, err
	}
	_, err = dao.UserDao.AddContacts(uid, group.Gid, dao.ContactsTypeGroup, "")
	if err != nil {
		// TODO undo
		return nil, err
	}
	GroupManager.AddGroup(g)
	GroupManager.PutMember(g.Gid, owner[0])
	return g, nil
}

func (m *groupApi) addGroupMember(gid int64, uid ...int64) ([]*dao.GroupMember, error) {

	g := GroupManager.GetGroup(gid)
	memberUid := make([]int64, 0, len(uid))
	for _, u := range uid {
		// member exist
		if !g.HasMember(u) {
			memberUid = append(memberUid, u)
		}
	}
	if len(memberUid) == 0 {
		return nil, errors.New("already added")
	}

	// TODO query user info and notify group members, optimize query time
	exist, err2 := dao.UserDao.HasUser(memberUid...)
	if err2 != nil {
		return nil, err2
	}
	if !exist {
		return nil, errors.New("user does not exist")
	}

	members, err := dao.GroupDao.AddMember(gid, dao.GroupMemberUser, memberUid...)
	if err != nil {
		return nil, err
	}
	return members, nil
}
