package im

import (
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
	ret.groups = NewGroupMap()
	return ret
}

func (m *groupManager) GetGroup(gid uint64) *Group {
	defer m.LockUtilReturn()()
	g := m.groups.Get(gid)
	if g != nil {
		return g
	}
	name, member := dao.GroupDao.GetGroup(gid)
	NewGroup(gid, name, member)
	m.groups.Put(gid, g)
	return g
}

func (m *groupManager) DispatchMessage(c *Client, message *entity.Message) error {

	groupMsg := new(entity.GroupMessage)
	err := message.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message error", err)
		return err
	}

	group := m.GetGroup(groupMsg.Gid)
	return group.SendMessage(c.uid, message)
}

type GroupMap struct {
	*mutex
	groupsMap map[uint64]*Group
}

func NewGroupMap() *GroupMap {
	ret := new(GroupMap)
	ret.groupsMap = make(map[uint64]*Group)
	return ret
}

func (g *GroupMap) Size() int {
	return len(g.groupsMap)
}

func (g *GroupMap) Get(gid uint64) *Group {
	defer g.LockUtilReturn()()
	group, ok := g.groupsMap[gid]
	if ok {
		return group
	}
	return nil
}

func (g *GroupMap) Put(gid uint64, group *Group) {
	defer g.LockUtilReturn()()
	g.groupsMap[gid] = group
}

func (g *GroupMap) Delete(gid uint64) {
	defer g.LockUtilReturn()()
	delete(g.groupsMap, gid)
}
