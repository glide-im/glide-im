package im

import (
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

func (m *groupManager) GetGroup(gid int64) *Group {
	defer m.LockUtilReturn()()
	g := m.groups.Get(gid)
	if g != nil {
		return g
	}
	//group, err := dao.GroupDao.GetGroup(gid)
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
