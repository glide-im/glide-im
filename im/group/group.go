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
	messages chan *message.DownGroupMessage
	// notify 群通知队列
	notify chan *message.GroupNotify

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
	ret.messages = make(chan *message.DownGroupMessage, 100)
	ret.notify = make(chan *message.GroupNotify, 10)
	ret.checkActive = tw.After(time.Minute * 30)
	ret.msgSequence = 1
	ret.gid = gid
	return ret
}

func (g *Group) EnqueueNotify(msg *message.GroupNotify) error {
	select {
	case g.notify <- msg:
		atomic.AddInt32(&g.queued, 1)
	default:
		return errors.New("notify message queue is full")
	}
	g.checkMsgQueue()
	return nil
}

func (g *Group) EnqueueMessage(msg *message.UpChatMessage) error {

	g.mu.Lock()
	mf, exist := g.members[msg.From_]
	g.mu.Unlock()

	if !exist {
		return errors.New("not a group member")
	}
	if mf.muted {
		return errors.New("a muted group member send message")
	}
	dMsg := &message.DownGroupMessage{
		Mid:     msg.Mid,
		MsgSeq:  atomic.AddInt64(&g.msgSequence, 1),
		From:    msg.From_,
		Type:    msg.Type,
		Content: msg.Content,
		SendAt:  msg.CTime,
	}

	select {
	case g.messages <- dMsg:
		atomic.AddInt32(&g.queued, 1)
	default:
		return errors.New("too many messages,the group message queue is full")
	}
	g.checkMsgQueue()
	return nil
}

func (g *Group) checkMsgQueue() {
	if atomic.LoadInt32(&g.queued) > 0 {
		return
	}
	err := queueExec.Submit(
		func() {
			logger.D("run a message queue reader goroutine")
			for {
				select {
				case m := <-g.notify:
					atomic.StoreInt32(&g.queued, -1)
					switch m.Type {
					case 1:
					case 2:
					case 3:
					}
					// 优先派送群通知消息
					continue
				case <-g.checkActive.C:
					if g.lastMsgAt.Add(time.Minute*30).After(time.Now()) && atomic.LoadInt32(&g.queued) == 0 {
						// 超过三十分钟没有发消息了, 停止消息下行任务
						goto REST
					} else {
						g.checkActive = tw.After(time.Minute * 30)
					}
				case m := <-g.messages:
					atomic.StoreInt32(&g.queued, -1)
					g.lastMsgAt = time.Now()
					g.SendMessage(m.From, message.NewMessage(-1, message.ActionGroupMessage, m))
				}
			}
		REST:
			logger.D("message queue read goroutine exit")
		},
	)
	if err == ants.ErrPoolOverload {
		logger.E("group message queue handle goroutine pool is overload")
	}
}

func (g *Group) SendMessage(from int64, message *message.Message) {
	logger.D("Group.SendMessage: %s", message)
	g.mu.Lock()
	for uid, mf := range g.members {
		if !mf.online || uid == from {
			continue
		}
		client.EnqueueMessage(uid, message)
	}
	g.mu.Unlock()
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
