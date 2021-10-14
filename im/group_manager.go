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

func (m *groupManager) PutMember(gid int64, mb map[int64]int32) {
	for k, v := range mb {
		m.getGroup(gid).PutMember(k, v)
	}
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

func (m *groupManager) AddGroup(gid int64) {
	// TODO
}

func (m *groupManager) RemoveGroup(gid int64) {
	// TODO
}

func (m *groupManager) ChangeStatus(gid int64, status int64) {

}

func (m *groupManager) GetGroup1(gid int64) *dao.Group {

	g := m.groups.Get(gid)
	if g != nil {
		return nil
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

	g = group.NewGroup(dbGroup)
	for _, member := range members {
		g.PutMember(member.Uid, member.Type)
	}
	m.groups.Put(gid, g)
	return dbGroup
}

func (m *groupManager) DispatchNotifyMessage(gid int64, message *message.Message) {
	g := m.getGroup(gid)
	if g != nil {
		g.SendMessage(message)
	}
}

func (m *groupManager) DispatchMessage(gid int64, msg *message.Message) {
	logger.D("GroupManager.HandleMessage: %s", msg)

	groupMsg := new(client.GroupMessage)
	err := msg.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message error", err)
		return
	}

	g := m.getGroup(groupMsg.TargetId)

	if g == nil {
		logger.E("dispatch group message", "group not exist")
		return
	}

	g.EnqueueMessage(groupMsg.Sender, groupMsg)
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
