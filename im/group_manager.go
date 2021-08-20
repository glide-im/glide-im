package im

import (
	"errors"
	"go_im/im/comm"
	"go_im/im/dao"
	"go_im/im/entity"
)

var GroupManager = NewGroupManager()

type groupManager struct {
	*comm.Mutex
	groups *groupMap
}

func NewGroupManager() *groupManager {
	ret := new(groupManager)
	ret.Mutex = new(comm.Mutex)
	ret.groups = NewGroupMap()
	return ret
}

func (m *groupManager) PutMember(gid int64, mb *dao.GroupMember) {
	m.GetGroup(gid).PutMember(mb)
}

func (m *groupManager) UnsubscribeGroup(uid int64, gid int64) {
	m.GetGroup(gid).RemoveMember(uid)
}

func (m *groupManager) RemoveMember(gid int64, uid ...int64) error {
	group := m.GetGroup(gid)
	if group == nil {
		return errors.New("unknown group")
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

func (m *groupManager) AddGroup(group *Group) {
	m.groups.Put(group.Gid, group)
}

func (m *groupManager) GetGroup(gid int64) *Group {

	g := m.groups.Get(gid)
	if g != nil {
		return g
	}

	group, err := dao.GroupDao.GetGroup(gid)
	if err != nil {
		comm.Slog.E("GroupManager.GetGroup", "load group", gid, err)
		return nil
	}

	members, err := dao.GroupDao.GetMembers(gid)
	if err != nil {
		comm.Slog.E("GroupManager.GetGroup", "load members", gid, err)
		return nil
	}

	chat, err := dao.ChatDao.GetChat(gid, 2)
	if err != nil {
		comm.Slog.E("GroupManager.GetGroup", "load chat", gid, err)
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
	comm.Slog.D("GroupManager.DispatchMessage: %s", message)

	groupMsg := new(entity.GroupMessage)
	err := message.DeserializeData(groupMsg)
	if err != nil {
		comm.Slog.E("dispatch group message error", err)
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

	resp := entity.NewMessage(-1, entity.ActionChatMessage, msg)

	group.SendMessage(uid, resp)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type groupMap struct {
	*comm.Mutex
	groupsMap map[int64]*Group
}

func NewGroupMap() *groupMap {
	ret := new(groupMap)
	ret.Mutex = new(comm.Mutex)
	ret.groupsMap = make(map[int64]*Group)
	return ret
}

func (g *groupMap) Size() int {
	return len(g.groupsMap)
}

func (g *groupMap) Get(gid int64) *Group {
	defer g.LockUtilReturn()()
	group, ok := g.groupsMap[gid]
	if ok {
		return group
	}
	return nil
}

func (g *groupMap) Put(gid int64, group *Group) {
	defer g.LockUtilReturn()()
	g.groupsMap[gid] = group
}

func (g *groupMap) Delete(gid int64) {
	defer g.LockUtilReturn()()
	delete(g.groupsMap, gid)
}
