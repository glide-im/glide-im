package group

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/pkg/timingwheel"
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

type memberInfo struct {
	online    bool
	muted     bool
	admin     bool
	deletedAt int64
}

func newMemberInfo() memberInfo {
	return memberInfo{
		online:    false,
		muted:     false,
		admin:     false,
		deletedAt: 0,
	}
}

var tw = timingwheel.NewTimingWheel(time.Second, 3, 20)
var queueExec *ants.Pool

func init() {
	var e error
	queueExec, e = ants.NewPool(200000,
		ants.WithNonblocking(true),
		ants.WithPreAlloc(true),
		ants.WithPanicHandler(onQueueExecutorPanic),
	)
	if e != nil {
		panic(e)
	}
}

func onQueueExecutorPanic(i interface{}) {
	logger.E("message queue goroutine pool handle message queue panic %v", i)
}

type Group struct {
	gid int64

	msgSequence int64
	startup     string

	mute      bool
	dissolved bool

	// messages 群消息队列
	messages chan *message.UpChatMessage
	// notify 群通知队列
	notify chan *message.DownGroupMessage

	queued int32

	// checkActive 定时检查群活跃情况
	checkActive *timingwheel.Task

	lastMsgAt time.Time
	mu        *comm.Mutex
	members   map[int64]memberInfo
}

func newGroup(gid int64) *Group {
	ret := new(Group)
	ret.mu = comm.NewMutex()
	ret.members = map[int64]memberInfo{}
	ret.startup = strconv.FormatInt(time.Now().Unix(), 10)
	ret.messages = make(chan *message.UpChatMessage, 100)
	ret.checkActive = tw.After(time.Minute * 30)
	ret.msgSequence = 1
	ret.gid = gid
	return ret
}

func (g *Group) EnqueueMessage(msg *message.UpChatMessage) {

	atomic.AddInt32(&g.queued, 1)
	if atomic.LoadInt32(&g.queued) > 1 {
		select {
		case g.messages <- msg:
		default:
			logger.E("too many messages,the group message queue is full")
		}
		return
	}

	err := queueExec.Submit(
		func() {
			for {
				select {
				case <-g.checkActive.C:
					if g.lastMsgAt.Add(time.Minute*30).After(time.Now()) && atomic.LoadInt32(&g.queued) == 0 {
						// 超过三十分钟没有发消息了, 停止消息下行任务
						goto REST
					} else {
						g.checkActive = tw.After(time.Minute * 30)
					}
				case m := <-g.messages:
					g.handleMessage(m)
					atomic.StoreInt32(&g.queued, -1)
					g.lastMsgAt = time.Now()
				}
			}
		REST:
		})
	if err == ants.ErrPoolOverload {
		logger.E("group message queue handle goroutine pool is overload")
	}
}

func (g *Group) SendMessage(message *message.Message) {
	logger.D("Group.SendMessage: %s", message)

	for uid, mf := range g.members {
		if !mf.online {
			continue
		}
		client.EnqueueMessage(uid, message)
	}
}

func (g *Group) updateMember(u MemberUpdate) error {
	defer g.mu.LockUtilReturn()()
	mf, ok := g.members[u.Uid]
	if !ok && u.Flag != FlagMemberAdd {
		return errors.New("member not exist")
	}
	switch u.Flag {
	case FlagMemberDel:
		mf.deletedAt = time.Now().Unix()
	case FlagMemberAdd:
		g.members[u.Uid] = newMemberInfo()
	case FlagMemberMuted:
		mf.muted = true
	case FlagMemberOnline:
		mf.online = true
	case FlagMemberOffline:
		mf.online = false
	case FlagMemberCancelAdmin:
		mf.admin = false
	case FlagMemberSetAdmin:
		mf.admin = true
	}
	return nil
}

func (g *Group) PutMember(member int64, s memberInfo) {
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

func (g *Group) handleMessage(msg *message.UpChatMessage) {
	g.mu.Lock()
	mf, exist := g.members[msg.From_]
	g.mu.Unlock()

	if !exist {
		logger.W("a non-group member send message")
		return
	}
	if mf.muted {
		logger.W("a muted group member send message")
		return
	}
	seq := atomic.LoadInt64(&g.msgSequence)
	rMsg := &message.DownGroupMessage{
		Mid:     msg.Mid,
		MsgSeq:  seq,
		From:    msg.From_,
		Type:    msg.Type,
		Content: msg.Content,
		SendAt:  msg.CTime,
	}

	resp := message.NewMessage(-1, message.ActionGroupMessage, rMsg)

	g.SendMessage(resp)
	atomic.AddInt64(&g.msgSequence, 1)
}
