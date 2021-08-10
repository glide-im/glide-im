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

func (m *groupManager) AddGroup(group *Group) {
	defer m.LockUtilReturn()()

	m.groups.Put(group.Gid, group)
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

func (m *groupManager) DispatchMessage(c *Client, message *entity.Message) error {
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
	chatMsg, err := dao.ChatDao.NewChatMessage(groupMsg.Cid, c.uid, groupMsg.Message, groupMsg.MessageType)
	if err != nil {
		return err
	}

	msg := entity.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         groupMsg.Cid,
		UcId:        groupMsg.UcId,
		Sender:      c.uid,
		MessageType: 1,
		Message:     groupMsg.Message,
		SendAt:      groupMsg.SendAt,
	}

	resp := entity.NewMessage2(-1, entity.ActionChatMessage, msg)

	group.SendMessage(c.uid, resp)
	return nil
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
