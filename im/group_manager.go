package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

var GroupManager = NewGroupManager()

type groupManager struct {
	*mutex
	groups *GroupMap
}

func NewGroupManager() *groupManager {
	ret := new(groupManager)
	ret.mutex = new(mutex)
	ret.groups = NewGroupMap()
	return ret
}

func (m *groupManager) CreateGroup(name string, uid int64) (*Group, error) {

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
	g.PutMember(owner[0])

	defer m.LockUtilReturn()()
	m.groups.Put(g.Gid, g)

	return g, nil
}

func (m groupManager) AddGroupMember(gid int64, uid ...int64) ([]*dao.GroupMember, error) {

	g := m.GetGroup(gid)
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

func (m *groupManager) RemoveMember(gid int64, uid ...int64) error {
	group := m.GetGroup(gid)
	if group == nil {
		return errors.New("")
	}
	for _, id := range uid {
		if group.HasMember(id) {
			group.RemoveMember(id)
		}
	}
	return nil
}

func (m *groupManager) GetMembers(gid int64) ([]*dao.GroupMember, error) {
	group := m.GetGroup(gid)
	if group == nil {
		return []*dao.GroupMember{}, nil
	}
	return group.GetMembers(), nil
}

func (m *groupManager) GetGroup(gid int64) *Group {
	defer m.LockUtilReturn()()

	g := m.groups.Get(gid)
	if g != nil {
		return g
	}

	group, err := dao.GroupDao.GetGroup(gid)
	if err != nil {
		logger.E("GroupManager.GetGroup", "load group", gid, err)
		return nil
	}

	members, err := dao.GroupDao.GetMembers(gid)
	if err != nil {
		logger.E("GroupManager.GetGroup", "load members", gid, err)
		return nil
	}

	chat, err := dao.ChatDao.GetChat(gid, 2)
	if err != nil {
		logger.E("GroupManager.GetGroup", "load chat", gid, err)
		return nil
	}
	g = NewGroup(gid, group, chat.Cid, members)
	m.groups.Put(gid, g)
	return g
}

func (m *groupManager) DispatchNotifyMessage(uid int64, gid int64, message *entity.Message) {
	group := m.GetGroup(gid)
	if group != nil {
		group.SendMessage(uid, message)
	}
}

func (m *groupManager) DispatchMessage(uid int64, message *entity.Message) error {
	logger.D("GroupManager.DispatchMessage: %s", message)

	groupMsg := new(entity.GroupMessage)
	err := message.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message error", err)
		return err
	}

	group := m.GetGroup(groupMsg.TargetId)

	if group == nil {
		return errors.New("group not exist")
	}

	//
	chatMsg, err := dao.ChatDao.NewChatMessage(groupMsg.Cid, uid, groupMsg.Message, groupMsg.MessageType)
	if err != nil {
		return err
	}

	msg := entity.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         groupMsg.Cid,
		UcId:        groupMsg.UcId,
		Sender:      uid,
		MessageType: 1,
		Message:     groupMsg.Message,
		SendAt:      groupMsg.SendAt,
	}

	resp := entity.NewMessage2(-1, entity.ActionChatMessage, msg)

	group.SendMessage(uid, resp)
	return nil
}

func (m *groupManager) SubscribeGroup(gid int64, mb *dao.GroupMember) {
	m.GetGroup(gid).PutMember(mb)
}

func (m *groupManager) UnsubscribeGroup(uid int64, gid int64) {
	m.GetGroup(gid).RemoveMember(uid)
}

////////////////////////////////////////////////////////////////////////////////

type GroupMap struct {
	*mutex
	groupsMap map[int64]*Group
}

func NewGroupMap() *GroupMap {
	ret := new(GroupMap)
	ret.mutex = new(mutex)
	ret.groupsMap = make(map[int64]*Group)
	return ret
}

func (g *GroupMap) Size() int {
	return len(g.groupsMap)
}

func (g *GroupMap) Get(gid int64) *Group {
	defer g.LockUtilReturn()()
	group, ok := g.groupsMap[gid]
	if ok {
		return group
	}
	return nil
}

func (g *GroupMap) Put(gid int64, group *Group) {
	defer g.LockUtilReturn()()
	g.groupsMap[gid] = group
}

func (g *GroupMap) Delete(gid int64) {
	defer g.LockUtilReturn()()
	delete(g.groupsMap, gid)
}
