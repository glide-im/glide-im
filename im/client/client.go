package client

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/conn"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
	"go_im/pkg/timingwheel"
	"sync/atomic"
	"time"
)

var tw = timingwheel.NewTimingWheel(time.Millisecond*500, 3, 20)

// HeartbeatDuration 心跳间隔
const HeartbeatDuration = time.Second * 5

var pool *ants.Pool

func init() {
	var err error
	pool, err = ants.NewPool(50_0000,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(func(i interface{}) {
			logger.E("")
		}),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		panic(err)
	}
}

const (
	_ = iota
	stateRunning
	stateClosing
	stateClosed
)

type Info struct {
	ID           int64
	AliveAt      int64
	ConnectionAt int64
	Device       int64
}

// IClient 表示一个客户端, 用于管理连接状态, 连接 id, 消息收发
type IClient interface {

	// SetID 设置该客户端标识 ID
	SetID(id int64, device int64)

	// Closed 返回该客户端连接是否已关闭
	Closed() bool

	// EnqueueMessage 将消息放入到客户端消息队列中
	EnqueueMessage(message *message.Message) error

	// Exit 退出客户端, 关闭连接等
	Exit()

	// Run 开始收发消息客户端连接的消息
	Run()

	GetInfo() Info
}

// Client represent a user conn conn
type Client struct {
	conn conn.Connection

	// id 唯一标识
	id int64
	// device 设备标识
	device    int64
	connectAt time.Time
	// state client 状态
	state int32

	// queuedMessage messages in the queue
	queuedMessage int64
	// messages 带缓冲的下行消息管道, 缓冲大小40
	messages chan *message.Message
	// rCloseCh 关闭或写入则停止读
	rCloseCh   chan struct{}
	readClosed int32

	// hbR 心跳倒计时
	hbR    *timingwheel.Task
	hbLost int

	hbW *timingwheel.Task

	// seq 服务器下行消息递增序列号
	seq int64
}

func newClient(conn conn.Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.state = stateRunning
	// 大小为 40 的缓冲管道, 防止短时间消息过多如果网络连接 output 不及时会造成程序阻塞, 可以适当调整
	client.messages = make(chan *message.Message, 60)
	client.connectAt = time.Now()
	client.rCloseCh = make(chan struct{})
	client.seq = 0
	client.hbR = tw.After(HeartbeatDuration)
	client.hbW = tw.After(HeartbeatDuration)
	return client
}

func (c *Client) GetInfo() Info {
	return Info{
		ID:           c.id,
		AliveAt:      0,
		ConnectionAt: c.connectAt.Unix(),
		Device:       c.device,
	}
}

// SetID 设置 id 标识及设备标识
func (c *Client) SetID(id int64, device int64) {
	//logger.D("set client id, origin: id=%d, device=%d, new: id=%d, device=%d", c.id, c.device, id, device)
	atomic.StoreInt64(&c.id, id)
	atomic.StoreInt64(&c.device, device)
}

func (c *Client) Closed() bool {
	return atomic.LoadInt32(&c.state) != stateRunning
}

// EnqueueMessage 放入下行消息队列
func (c *Client) EnqueueMessage(message *message.Message) error {
	atomic.AddInt64(&c.queuedMessage, 1)
	err := pool.Submit(func() {
		defer func() {
			atomic.AddInt64(&c.queuedMessage, -1)
			e := recover()
			if e != nil {
				logger.E("%v", e)
			}
		}()
		s := atomic.LoadInt32(&c.state)
		if s == stateClosed {
			logger.D("client has closed, enqueue message failed")
			return
		}
		logger.I("EnqueueMessage(id=%d, %s): %v", c.id, message.GetAction(), message)
		if message.GetSeq() < 0 {
			// 服务端主动发送的消息使用服务端的序列号
			message.SetSeq(c.getNextSeq())
		}
		select {
		case c.messages <- message:
		default:
			atomic.AddInt64(&c.queuedMessage, -1)
			// 消息 chan 缓冲溢出, 这条消息将被丢弃
			logger.E("message chan is full, id=%d", c.id)
		}
	})
	if err != nil {
		logger.E("message not enqueue:%v", err)
	}

	return nil
}

// readMessage 开始从 Connection 中读取消息
func (c *Client) readMessage() {
	readChan, done := messageReader.ReadCh(c.conn)

	defer func() {
		err := recover()
		if err != nil {
			logger.E("read message error", err)
		}
	}()

	atomic.StoreInt32(&c.readClosed, 0)
	for {
		select {
		case <-c.rCloseCh:
			close(c.rCloseCh)
			goto STOP
		case <-c.hbR.C:
			c.hbLost++
			if c.hbLost > 3 {
				logger.D("heartbeat timout, id=%d, device=%d", c.id, c.device)
				goto STOP
			}
			// reset client heartbeat
			c.hbR.Cancel()
			c.hbR = tw.After(HeartbeatDuration)
			c.EnqueueMessage(message.NewMessage(0, message.ActionHeartbeat, ""))
		case msg := <-readChan:
			if msg.err != nil {
				if c.Closed() || c.handleError(msg.err) {
					// 连接断开或致命错误中断读消息
					goto STOP
				}
				continue
			}
			c.hbLost = 0
			c.hbR.Cancel()
			c.hbR = tw.After(HeartbeatDuration)
			id, device := c.getID()
			// 统一处理消息函数
			_ = messageHandleFunc(id, device, msg.m)
			msg.Recycle()
		}
	}
STOP:
	c.hbR.Cancel()
	atomic.StoreInt32(&c.readClosed, 1)
	close(done)
	id, device := c.getID()
	logger.D("client read closed, id=%d, device=%d", id, device)
}

// writeMessage 开始向 Connection 中写入消息队列中的消息
func (c *Client) writeMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("Client write message error: %v", err)
		}
	}()

	for {
		select {
		case <-c.hbW.C:
			if c.Closed() {
				logger.D("read closed, down msg queue timeout, close write now, uid=%d", c.id)
				goto STOP
			}
			c.EnqueueMessage(message.NewMessage(c.getNextSeq(), message.ActionHeartbeat, struct{}{}))
			c.hbW.Cancel()
			c.hbW = tw.After(HeartbeatDuration)
		case m := <-c.messages:
			b, err := codec.Encode(m)
			if err != nil {
				logger.E("serialize output message", err)
				continue
			}
			err = c.conn.Write(b)
			atomic.AddInt64(&c.queuedMessage, -1)

			c.hbW.Cancel()
			c.hbW = tw.After(HeartbeatDuration)
			if err != nil {
				if c.Closed() || c.handleError(err) {
					// 连接断开或致命错误中断写消息
					goto STOP
				}
			} else {
				statistics.SMsgOutput()
			}
		}
	}
STOP:
	c.Exit()
	atomic.StoreInt32(&c.state, stateClosed)
	close(c.messages)
	_ = c.conn.Close()
	logger.D("client write closed, uid=%d", c.id)
}

// handleError 处理上下行消息过程中的错误, 如果是致命错误, 则返回 true
func (c *Client) handleError(err error) bool {
	statistics.SError(err)
	if conn.ErrClosed != err {
		logger.E("handle message error: %s", err.Error())
	}
	if !uid.IsTempId(atomic.LoadInt64(&c.id)) {
		err = Logout(atomic.LoadInt64(&c.id), c.device)
		if err != nil {
			logger.E("%v", err)
		}
	}
	return true
}

// Exit 退出客户端
func (c *Client) Exit() {
	s := atomic.LoadInt32(&c.state)
	if s == stateClosed || s == stateClosing {
		return
	}
	atomic.StoreInt32(&c.state, stateClosing)

	if atomic.LoadInt32(&c.readClosed) != 1 {
		c.rCloseCh <- struct{}{}
	}
}

func (c *Client) close() {

}

func (c *Client) getID() (int64, int64) {
	return atomic.LoadInt64(&c.id), atomic.LoadInt64(&c.device)
}

// getNextSeq 获取下一个下行消息序列号 sequence
func (c *Client) getNextSeq() int64 {
	return atomic.AddInt64(&c.seq, 1)
}

func (c *Client) Run() {
	logger.I(">>>> client %s running, id=%d", c.conn.GetConnInfo().Addr, c.id)
	go c.readMessage()
	go c.writeMessage()
}
