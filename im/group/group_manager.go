package group

import (
	"errors"
	"go_im/im/comm"
	"go_im/im/message"
	"time"
)

// Manager 群相关操作入口
var Manager IGroupManager = NewDefaultManager()

type EnqueueMessageInterface interface {
	EnqueueMessage(uid int64, device int64, message *message.Message)
}

var EnqueueMessage EnqueueMessageInterface

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

	Extra interface{}
}

type Update struct {
	Flag int64

	Extra interface{}
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

// TODO 2021-11-20 大群小群优化

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
		m.mu.Lock()
		m.groups[gid] = newGroup(gid)
		m.mu.Unlock()
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
		tw.After(time.Second * 10).Callback(func() {
			m.mu.Lock()
			delete(m.groups, gid)
			m.mu.Unlock()
		})
	}
	return nil
}

func (m *DefaultManager) DispatchNotifyMessage(gid int64, msg *message.Message) error {
	m.mu.Lock()
	g := m.groups[gid]
	m.mu.Unlock()
	return g.EnqueueNotify(&message.GroupNotify{})
}

func (m *DefaultManager) DispatchMessage(gid int64, msg *message.UpChatMessage) error {
	//logger.D("GroupManager.HandleMessage: %v", msg)
	m.mu.Lock()
	g, ok := m.groups[gid]
	m.mu.Unlock()
	if !ok {
		return errors.New("group not exist")
	}
	if g.mute {
		return errors.New("group is muted")
	}
	if g.dissolved {
		return errors.New("group is dissolved")
	}
	return g.EnqueueMessage(msg)
}
