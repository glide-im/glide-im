package im

import (
	"errors"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type groupManager struct {
	mu     *comm.Mutex
	groups map[int64]*group.Group
}

func NewGroupManager() *groupManager {
	ret := new(groupManager)
	ret.mu = comm.NewMutex()
	ret.groups = map[int64]*group.Group{}
	return ret
}

func (m *groupManager) init() {
	allGroup := group.LoadAllGroup()
	for gid, g := range allGroup {
		m.groups[gid] = g
	}
}

func (m *groupManager) PutMember(gid int64, mb map[int64]int32) {
	g := m.groups[gid]
	for k, v := range mb {
		g.PutMember(k, v)
	}
}

func (m *groupManager) RemoveMember(gid int64, uid ...int64) error {
	g := m.groups[gid]
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
	g, err := group.LoadGroup(gid)
	if err != nil {
		return
	}
	m.groups[gid] = g
}

func (m *groupManager) RemoveGroup(gid int64) {
	g := m.groups[gid]
	if g != nil {

	}
}

func (m *groupManager) ChangeStatus(gid int64, status int64) {

}

func (m *groupManager) DispatchNotifyMessage(gid int64, message *message.Message) {
	g := m.groups[gid]
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

	g := m.groups[gid]

	if g == nil {
		logger.E("dispatch group message", "group not exist")
		return
	}

	g.EnqueueMessage(groupMsg.Sender, groupMsg)
}
