package im

import "go_im/im/entity"

var GroupManager = NewGroupManager()

type groupManager struct {
	groups *GroupMap
}

func NewGroupManager() *groupManager {
	ret := new(groupManager)
	ret.groups = NewGroupMap()
	return ret
}

func (m *groupManager) DispatchMessage(c *Client, message *entity.Message) {

	groupMsg := new(entity.GroupMessageEntity)
	err := message.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message error", err)
		return
	}

	group := m.groups.Get(groupMsg.Gid)
	group.SendMessage(message)
}
