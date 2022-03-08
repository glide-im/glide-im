package group

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/pkg/timingwheel"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type memberInfo struct {
	online    bool
	muted     bool
	admin     bool
	deletedAt int64
}

func newMemberInfo() *memberInfo {
	return &memberInfo{
		online:    false,
		muted:     false,
		admin:     false,
		deletedAt: 0,
	}
}

var tw = timingwheel.NewTimingWheel(time.Second, 3, 20)
var queueExec *ants.Pool

const messageQueueSleep = time.Second * 10

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
	messages chan *message.ChatMessage
	// notify 群通知队列
	notify chan *message.GroupNotify

	queueRunning int32
	queued       int32

	// checkActive 定时检查群活跃情况
	checkActive *timingwheel.Task

	lastMsgAt time.Time
	mu        *sync.Mutex
	members   map[int64]*memberInfo
}

func newGroup(gid int64) *Group {
	ret := new(Group)
	ret.mu = &sync.Mutex{}
	ret.members = map[int64]*memberInfo{}
	ret.startup = strconv.FormatInt(time.Now().Unix(), 10)
	ret.messages = make(chan *message.ChatMessage, 100)
	ret.notify = make(chan *message.GroupNotify, 10)
	ret.checkActive = tw.After(messageQueueSleep)
	ret.queueRunning = 0
	ret.msgSequence = 1
	ret.gid = gid
	return ret
}

func (g *Group) EnqueueNotify(msg *message.GroupNotify) error {
	seq := atomic.AddInt64(&g.msgSequence, 1)
	msg.Seq = seq
	select {
	case g.notify <- msg:
		atomic.AddInt32(&g.queued, 1)
	default:
		return errors.New("notify message queue is full")
	}
	return g.checkMsgQueue()
}

func (g *Group) EnqueueMessage(msg *message.ChatMessage, recall bool) (int64, error) {

	g.mu.Lock()
	mf, exist := g.members[msg.From]
	g.mu.Unlock()

	if !exist {
		return 0, errors.New("not a group member")
	}
	if mf.muted && !recall {
		return 0, errors.New("a muted group member send message")
	}

	if recall {
		r := &message.Recall{}
		err := message.UnmarshallJson(msg.Content, r)
		if err != nil {
			return 0, err
		}
		if r.RecallBy != msg.From && !mf.admin {
			return 0, errors.New("illegal operation")
		}
		err = msgdao.GroupMsgDaoImpl.UpdateGroupMessageRecall(g.gid, r.Mid, msgdao.ChatMessageStatusRecalled, r.RecallBy)
		if err != nil {
			return 0, err
		}
	}
	now := time.Now().Unix()
	seq := atomic.AddInt64(&g.msgSequence, 1)
	err := msgdao.AddGroupMessage(&msgdao.GroupMessage{
		MID:     msg.Mid,
		Seq:     seq,
		To:      g.gid,
		From:    msg.From,
		Type:    msg.Type,
		SendAt:  now,
		Content: msg.Content,
	})
	if err != nil {
		atomic.AddInt64(&g.msgSequence, -1)
		return 0, err
	}
	err = msgdao.UpdateGroupMessageState(g.gid, msg.Mid, time.Now().Unix(), seq)
	if err != nil {
		logger.E("Group.EnqueueMessage update group message state error, %v", err)
		return 0, err
	}

	dMsg := msg
	msg.Seq = seq
	msg.To = g.gid
	msg.SendAt = now
	if err != nil {
		return 0, err
	}

	select {
	case g.messages <- dMsg:
		atomic.AddInt32(&g.queued, 1)
	default:
		return 0, errors.New("too many messages,the group message queue is full")
	}
	if err := g.checkMsgQueue(); err != nil {
		if err == ants.ErrPoolOverload {
			logger.E("group message queue handle goroutine pool is overload")
		}
		return 0, err
	}
	return seq, nil
}

func (g *Group) checkMsgQueue() error {
	if atomic.LoadInt32(&g.queueRunning) == 1 {
		return nil
	}
	err := queueExec.Submit(
		func() {
			atomic.StoreInt32(&g.queueRunning, 1)
			logger.D("run a message queue reader goroutine")
			g.checkActive = tw.After(messageQueueSleep)
			for {
				select {
				case m := <-g.notify:
					g.lastMsgAt = time.Now()
					atomic.AddInt32(&g.queued, -1)
					switch m.Type {
					default:
						g.SendMessage(0, message.NewMessage(0, message.ActionNotifyGroup, m))
					}
					// 优先派送群通知消息
					continue
				case <-g.checkActive.C:
					g.checkActive.Cancel()
					if g.lastMsgAt.Add(messageQueueSleep).Before(time.Now()) {
						q := atomic.LoadInt32(&g.queued)
						if q != 0 {
							logger.W("group message queue blocked, size=" + strconv.FormatInt(int64(q), 10))
							return
						}
						// 超过三十分钟没有发消息了, 停止消息下行任务
						goto REST
					} else {
						g.checkActive = tw.After(messageQueueSleep)
					}
				case m := <-g.messages:
					atomic.AddInt32(&g.queued, -1)
					g.lastMsgAt = time.Now()
					g.SendMessage(m.From, message.NewMessage(-1, message.ActionGroupMessage, m))
				}
			}
		REST:
			logger.D("message queue read goroutine exit")
			atomic.StoreInt32(&g.queueRunning, 0)
		},
	)
	if err != nil {
		atomic.StoreInt32(&g.queueRunning, 0)
	}
	return err
}

func (g *Group) SendMessage(from int64, message *message.Message) {
	logger.D("Group.SendMessage: %s", message)
	g.mu.Lock()
	for uid, mf := range g.members {
		if !mf.online || uid == from {
			continue
		}
		enqueueMessage(uid, 0, message)
	}
	g.mu.Unlock()
}

func (g *Group) updateMember(u MemberUpdate) error {
	g.mu.Lock()
	defer g.mu.Unlock()
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

func (g *Group) GetMember(id int64) *memberInfo {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.members[id]
}

func (g *Group) PutMember(member int64, s *memberInfo) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.members[member] = s
}

func (g *Group) RemoveMember(uid int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.members, uid)
}

func (g *Group) HasMember(uid int64) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, exist := g.members[uid]
	return exist
}
