package group

import (
	"errors"
	"go_im/im/comm"
	"go_im/im/message"
	"go_im/pkg/logger"
)

// Manager 群相关操作入口
var Manager IGroupManager = NewDefaultManager()

const (
	_ = iota
	FlagMemberAdd
	FlagMemberDel
	FlagMemberOnline
	FlagMemberOffline
	FlagMemberMuted
	FlagMemberSetAdmin
	FlagMemberCancelAdmin
)

const (
	_ = iota
	FlagGroupCreate
	FlagGroupDissolve
	FlagGroupMute
	FlagGroupCancelMute
)

type MemberUpdate struct {
	Uid  int64
	Flag int64
}

type Update struct {
	Flag int64
}

type IGroupManager interface {
	// UpdateMember 更新群成员
	UpdateMember(gid int64, update []MemberUpdate) error

	// UpdateGroup 更新群
	UpdateGroup(gid int64, update Update) error

	// DispatchNotifyMessage 发送通知消息
	DispatchNotifyMessage(gid int64, message *message.Message) error

	// DispatchMessage 发送聊天消息
	DispatchMessage(gid int64, message *message.UpChatMessage) error
}

type DefaultManager struct {
	mu     *comm.Mutex
	groups map[int64]*Group
}

func NewDefaultManager() *DefaultManager {
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

func (m *DefaultManager) UpdateMember(gid int64, update []MemberUpdate) error {
	for _, mbUpdate := range update {
		m.mu.Lock()
		g, ok := m.groups[gid]
		m.mu.Unlock()
		if !ok {
			return errors.New("group not exist")
		}
		err := g.updateMember(mbUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) UpdateGroup(gid int64, update Update) error {

	if update.Flag == FlagGroupCreate {
		m.groups[gid] = newGroup(gid)
		return nil
	}
	m.mu.Lock()
	g, ok := m.groups[gid]
	m.mu.Unlock()
	if !ok {
		return errors.New("group not exist")
	}
	switch update.Flag {
	case FlagGroupMute:
		g.mute = true
	case FlagGroupCancelMute:
		g.mute = false
	case FlagGroupDissolve:
		g.dissolved = true
	}
	return nil
}

func (m *DefaultManager) DispatchNotifyMessage(gid int64, message *message.Message) error {
	m.mu.Lock()
	g := m.groups[gid]
	m.mu.Unlock()
	if g != nil {
		g.SendMessage(message)
	}
	return nil
}

func (m *DefaultManager) DispatchMessage(gid int64, msg *message.UpChatMessage) error {
	logger.D("GroupManager.HandleMessage: %v", msg)
	m.mu.Lock()
	g, ok := m.groups[gid]
	m.mu.Unlock()
	if !ok {
		return nil
	}
	if g.mute {
		return nil
	}
	if g.dissolved {
		return nil
	}
	g.EnqueueMessage(msg)
	return nil
}
