package group

import (
	"errors"
	"go_im/im/client"
	"go_im/im/dao/groupdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strconv"
	"sync"
	"time"
)

// TODO 2021-11-20 大群小群优化

type DefaultManager struct {
	mu     *sync.Mutex
	groups map[int64]*Group
	h      MessageHandler
}

func NewDefaultManager() *DefaultManager {
	ret := new(DefaultManager)
	ret.mu = &sync.Mutex{}
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

func (m *DefaultManager) DispatchMessage(gid int64, action message.Action, msg *message.ChatMessage) error {
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
		ack := message.NewMessage(0, message.ActionAckNotify, message.NewAckMessage(msg.Mid, seq))
		client.EnqueueMessage(msg.From, ack)
	}

	return nil
}
