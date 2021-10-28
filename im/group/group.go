package group

import (
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	FlagShiftCanSend    = 1
	FlagShiftCanReceive = 2
	FlagShiftIsManager  = 3

	FlagDefault = 1 << FlagShiftCanSend
)

type Group struct {
	gid int64
	cid int64

	msgSequence int64
	startup     string

	mute bool

	mu      *comm.Mutex
	members map[int64]int32
}

func newGroup(gid int64, cid int64) *Group {
	ret := new(Group)
	ret.mu = comm.NewMutex()
	ret.members = map[int64]int32{}
	ret.startup = strconv.FormatInt(time.Now().Unix(), 10)
	ret.msgSequence = 1
	ret.gid = gid
	ret.cid = cid
	return ret
}

func (g *Group) EnqueueMessage(msg *message.GroupMessage) {

	flag, exist := g.members[msg.Sender]
	if !exist {
		logger.W("a non-group member send message")
		return
	}
	if flag&(1<<FlagShiftCanSend) == 0 {
		logger.W("a muted group member send message")
		return
	}

	seq := atomic.LoadInt64(&g.msgSequence)
	mid := dao.GetNextMessageId(g.cid)

	chatMessage := dao.ChatMessage{
		Mid:         mid,
		Cid:         g.cid,
		Seq:         seq,
		Sender:      msg.Sender,
		SendAt:      dao.Timestamp(time.Now()),
		Message:     msg.Message,
		MessageType: msg.MessageType,
		At:          "",
	}
	err := dao.ChatDao.NewGroupMessage(chatMessage)

	if err != nil {
		logger.E("dispatch group message", err)
		return
	}

	rMsg := message.ReceiverChatMessage{
		Mid:         mid,
		Cid:         g.cid,
		AlignTag:    g.startup,
		Seq:         seq,
		Sender:      msg.Sender,
		MessageType: msg.MessageType,
		Message:     msg.Message,
		SendAt:      msg.SendAt,
	}

	resp := message.NewMessage(-1, message.ActionGroupMessage, rMsg)

	g.SendMessage(resp)
	atomic.AddInt64(&g.msgSequence, 1)
}

func (g *Group) SendMessage(message *message.Message) {
	logger.D("Group.SendMessage: %s", message)

	for uid, flag := range g.members {
		if flag&(1<<FlagShiftCanReceive) == 1 {
			continue
		}
		client.EnqueueMessage(uid, message)
	}
}

func (g *Group) PutMember(member int64, s int32) {
	defer g.mu.LockUtilReturn()()
	g.members[member] = s
}

func (g *Group) RemoveMember(uid int64) {
	defer g.mu.LockUtilReturn()()
	delete(g.members, uid)
}

func (g *Group) HasMember(uid int64) bool {
	defer g.mu.LockUtilReturn()()
	_, exist := g.members[uid]
	return exist
}
