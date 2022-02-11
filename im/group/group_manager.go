package group

import (
	"errors"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/dao/groupdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strconv"
	"time"
)

// Manager 群相关操作入口
var Manager IGroupManager = NewDefaultManager()

type EnqueueMessageInterface interface {
	EnqueueMessage(uid int64, device int64, message *message.Message)
}

var EnqueueMessage EnqueueMessageInterface = client.Manager

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
	DispatchNotifyMessage(gid int64, message *message.GroupNotify) error

	// DispatchMessage 发送聊天消息
	DispatchMessage(gid int64, action message.Action, message *message.UpChatMessage) error
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
	groups, err := groupdao.Dao.GetAllGroup()
	if err != nil {
		logger.E("Init group error", err)
		return
	}
	for _, g := range groups {
		logger.D("load group %d", g.Gid)
		sGroup := newGroup(g.Gid)
		m.groups[g.Gid] = sGroup
		mbs, err := groupdao.Dao.GetMembers(g.Gid)
		if err != nil {
			logger.E("load group member error gid=%d %v", g.Gid, err)
			continue
		}
		for _, mb := range mbs {
			info := newMemberInfo()
			info.muted = mb.Flag == 1
			info.admin = mb.Type == 1
			info.online = true
			sGroup.PutMember(mb.Uid, info)
		}
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

func (m *DefaultManager) DispatchNotifyMessage(gid int64, msg *message.GroupNotify) error {
	m.mu.Lock()
	g := m.groups[gid]
	m.mu.Unlock()
	return g.EnqueueNotify(msg)
}

func (m *DefaultManager) DispatchMessage(gid int64, action message.Action, msg *message.UpChatMessage) error {
	//logger.D("GroupManager.HandleMessage: %v", msg)
	m.mu.Lock()
	g, ok := m.groups[gid]
	m.mu.Unlock()
	if !ok {
		return errors.New("group not exist gid=" + strconv.FormatInt(gid, 10))
	}
	if g.mute && action != message.ActionGroupMessageRecall {
		return errors.New("group is muted")
	}
	if g.dissolved {
		return errors.New("group is dissolved")
	}
	seq, err := g.EnqueueMessage(msg, action == message.ActionGroupMessageRecall)

	if err != nil {
		return err
	} else {
		// notify sender, group message send successful
		ack := message.NewMessage(0, message.ActionAckNotify, message.AckMessage{Mid: msg.Mid, Seq: seq})
		client.EnqueueMessage(msg.From, ack)
	}

	return nil
}
