package im

import (
	"errors"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

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

func (m *groupManager) GetGroupCid(gid int64) int64 {
	return m.getGroup(gid).Cid
}

func (m *groupManager) HasMember(gid int64, uid int64) bool {
	return m.getGroup(gid).HasMember(uid)
}

func (m *groupManager) PutMember(gid int64, mb *dao.GroupMember) {
	m.getGroup(gid).PutMember(mb)
}

func (m *groupManager) RemoveMember(gid int64, uid ...int64) error {
	g := m.getGroup(gid)
	if g == nil {
		return errors.New("unknown group")
	}
	for _, id := range uid {
		if g.HasMember(id) {
			g.RemoveMember(id)
		}
	}
	return nil
}

func (m *groupManager) UserOnline(uid, gid int64) {

}

func (m *groupManager) UserOffline(uid, gid int64) {

}

func (m *groupManager) GetMembers(gid int64) ([]*dao.GroupMember, error) {
	g := m.getGroup(gid)
	if g == nil {
		return []*dao.GroupMember{}, nil
	}
	return g.GetMembers(), nil
}

func (m *groupManager) AddGroup(g *dao.Group, cid int64, owner *dao.GroupMember) {
	gp := group.NewGroup(g.Gid, g, 1, []*dao.GroupMember{})
	m.groups.Put(g.Gid, gp)
}

func (m *groupManager) GetGroup(gid int64) *dao.Group {

	g := m.groups.Get(gid)
	if g != nil {
		return g.Group
	}

	dbGroup, err := dao.GroupDao.GetGroup(gid)
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
	g = group.NewGroup(gid, dbGroup, chat.Cid, members)
	m.groups.Put(gid, g)
	return g.Group
}

func (m *groupManager) DispatchNotifyMessage(uid int64, gid int64, message *message.Message) {
	g := m.getGroup(gid)
	if g != nil {
		g.SendMessage(uid, message)
	}
}

func (m *groupManager) DispatchMessage(uid int64, msg *message.Message) error {
	logger.D("GroupManager.DispatchMessage: %s", msg)

	groupMsg := new(client.GroupMessage)
	err := msg.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message error", err)
		return err
	}

	g := m.getGroup(groupMsg.TargetId)

	if g == nil {
		return errors.New("group not exist")
	}

	//
	chatMsg, err := dao.ChatDao.NewChatMessage(groupMsg.Cid, uid, groupMsg.Message, groupMsg.MessageType)
	if err != nil {
		return err
	}

	rMsg := client.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         groupMsg.Cid,
		UcId:        groupMsg.UcId,
		Sender:      uid,
		MessageType: 1,
		Message:     groupMsg.Message,
		SendAt:      groupMsg.SendAt,
	}

	resp := message.NewMessage(-1, message.ActionChatMessage, rMsg)

	g.SendMessage(uid, resp)
	return nil
}

func (m *groupManager) getGroup(gid int64) *group.Group {
	return m.groups.Get(gid)
}

////////////////////////////////////////////////////////////////////////////////

type groupMap struct {
	*comm.Mutex
	groupsMap map[int64]*group.Group
}

func NewGroupMap() *groupMap {
	ret := new(groupMap)
	ret.Mutex = new(comm.Mutex)
	ret.groupsMap = make(map[int64]*group.Group)
	return ret
}

func (g *groupMap) Size() int {
	return len(g.groupsMap)
}

func (g *groupMap) Get(gid int64) *group.Group {
	defer g.LockUtilReturn()()
	gp, ok := g.groupsMap[gid]
	if ok {
		return gp
	}
	return nil
}

func (g *groupMap) Put(gid int64, group *group.Group) {
	defer g.LockUtilReturn()()
	g.groupsMap[gid] = group
}

func (g *groupMap) Delete(gid int64) {
	defer g.LockUtilReturn()()
	delete(g.groupsMap, gid)
}
