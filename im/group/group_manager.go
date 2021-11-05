package group

import (
	"errors"
	"go_im/im/comm"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type DefaultManager struct {
	mu     *comm.Mutex
	groups map[int64]*Group
}

func NewGroupManager() *DefaultManager {
	ret := new(DefaultManager)
	ret.mu = comm.NewMutex()
	ret.groups = map[int64]*Group{}
	return ret
}

func (m *DefaultManager) Init() {
	allGroup := LoadAllGroup()
	for gid, g := range allGroup {
		m.groups[gid] = g
	}
}

func (m *DefaultManager) PutMember(gid int64, mb map[int64]int32) {
	g := m.groups[gid]
	for k := range mb {
		var flag int32 = FlagDefault
		g.PutMember(k, flag)
	}
}

func (m *DefaultManager) RemoveMember(gid int64, uid ...int64) error {
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

func (m *DefaultManager) AddGroup(gid int64) {
	g, err := LoadGroup(gid)
	if err != nil {
		return
	}
	m.groups[gid] = g
}

func (m *DefaultManager) RemoveGroup(gid int64) {
	g := m.groups[gid]
	if g != nil {

	}
}

func (m *DefaultManager) ChangeStatus(gid int64, status int64) {

}

func (m *DefaultManager) DispatchNotifyMessage(gid int64, message *message.Message) {
	g := m.groups[gid]
	if g != nil {
		g.SendMessage(message)
	}
}

func (m *DefaultManager) DispatchMessage(gid int64, msg *message.GroupMessage) {
	logger.D("GroupManager.HandleMessage: %v", msg)

	g := m.groups[gid]

	if g == nil {
		logger.E("dispatch group message", "group not exist")
		return
	}

	g.EnqueueMessage(msg)
}
