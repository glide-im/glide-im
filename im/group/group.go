package group

import (
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"time"
)

const (
	FlagShiftCanSend    = 1
	FlagShiftCanReceive = 2
	FlagShiftIsManager  = 3
)

type Group struct {
	Gid int64
	Cid int64

	nextMid int64
	group   *dao.Group
	mute    bool

	mu      *comm.Mutex
	members map[int64]int32
}

func NewGroup(group *dao.Group) *Group {
	ret := new(Group)
	ret.mu = comm.NewMutex()
	ret.members = map[int64]int32{}
	ret.Gid = group.Gid
	ret.Cid = group.ChatId
	ret.group = group
	chat, err := dao.ChatDao.GetChat(group.ChatId)
	if err != nil {
		logger.E("new group error, chat not exist", err)
		return nil
	}
	ret.nextMid = chat.CurrentMid + 1
	return ret
}

func (g *Group) EnqueueMessage(sender int64, msg *client.GroupMessage) {

	flag, exist := g.members[sender]
	if !exist {
		logger.W("a non-group member send message")
		return
	}
	if flag&(1<<FlagShiftCanSend) == 0 {
		logger.W("a muted group member send message")
		return
	}

	chatMessage := dao.ChatMessage{
		Mid:         g.nextMid,
		Cid:         g.Cid,
		Sender:      sender,
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

	rMsg := client.ReceiverChatMessage{
		Mid:         g.nextMid,
		Cid:         g.Cid,
		Sender:      sender,
		MessageType: msg.MessageType,
		Message:     msg.Message,
		SendAt:      msg.SendAt,
	}

	resp := message.NewMessage(-1, message.ActionChatMessage, rMsg)

	g.SendMessage(resp)

	g.nextMid = dao.GetNextMessageId(g.Cid)
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
