package client

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
	"go_im/pkg/timingwheel"
	"sync/atomic"
	"time"
)

// MessageHandler 用于处理客户端消息
type MessageHandler func(from int64, device int64, message *message.Message)

// MessageHandleFunc 所有客户端消息都传递到该函数处理
var MessageHandleFunc MessageHandler = nil

var tw = timingwheel.NewTimingWheel(time.Millisecond*500, 3, 20)

const (
	ExitCodeTTL        = 1
	ExitCodeBySrv      = 2
	ExitCodeLoginMutex = 3
	ExitCodeByUser     = 4
)

const HeartbeatDuration = time.Minute * 8

// IClient 表示一个客户端, 用于管理连接状态, 连接 id, 消息收发
type IClient interface {

	// SetID 设置该客户端标识 ID
	SetID(id int64, device int64)

	// Closed 返回该客户端连接是否已关闭
	Closed() bool

	// EnqueueMessage 将消息放入到客户端消息队列中
	EnqueueMessage(message *message.Message)

	// Exit 退出客户端, 关闭连接等
	Exit(code int64)

	// Run 开始收发消息客户端连接的消息
	Run()
}

// Client represent a user conn conn
type Client struct {
	conn conn.Connection

	id        int64
	device    int64
	connectAt time.Time
	closed    *comm.AtomicBool

	// buffered channel 40
	messages chan *message.Message

	hb *timingwheel.Task

	seq *comm.AtomicInt64
}

func NewClient(conn conn.Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = comm.NewAtomicBool(false)
	// 大小为 40 的缓冲管道, 防止短时间消息过多如果网络连接 output 不及时会造成程序阻塞, 可以适当调整
	client.messages = make(chan *message.Message, 40)
	client.connectAt = time.Now()
	client.seq = new(comm.AtomicInt64)
	client.hb = tw.After(HeartbeatDuration)
	return client
}

func (c *Client) SetID(id int64, device int64) {
	atomic.StoreInt64(&c.id, id)
	atomic.StoreInt64(&c.device, device)
	// TODO 恢复 SEQ 序列号
}

func (c *Client) Closed() bool {
	return c.closed.Get()
}

func (c *Client) EnqueueMessage(message *message.Message) {
	logger.I("EnqueueMessage(id=%d, %s): %v", c.id, message.Action, message)
	if c.closed.Get() {
		logger.W("connection closed, cannot enqueue message")
		return
	}
	if message.Seq <= 0 {
		// 服务端主动发送的消息使用服务端的序列号
		message.Seq = c.getNextSeq()
	}
	select {
	case c.messages <- message:
	default:
		// TODO 客户端弱网消息下行速度过慢导致缓冲溢出
		// 消息 chan 缓冲溢出, 这条消息将被丢弃
		logger.E("message chan is full, id=%d", c.id)
	}
}

// readMessage 开始从 Connection 中读取消息
func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.E("read message error", err)
		}
		Manager.ClientLogout(c.id, c.device)
	}()

	for {
		msg, err := messageReader.Read(c.conn)
		if err != nil {
			if c.Closed() || c.handleError(err) {
				// 连接断开或致命错误中断读消息
				break
			}
			continue
		}
		if msg.Action == message.ActionHeartbeat {
			c.hb.Cancel()
			c.hb = tw.After(HeartbeatDuration)
		} else {
			MessageHandleFunc(c.id, c.device, msg)
		}
	}
}

// writeMessage 开始向 Connection 中写入消息队列中的消息
func (c *Client) writeMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("Client write message error: %v", err)
		}
		Manager.ClientLogout(c.id, c.device)
	}()

	for msg := range c.messages {
		b, err := msg.Serialize()
		if err != nil {
			logger.E("serialize output message", err)
			continue
		}
		err = c.conn.Write(b)
		if err != nil {
			if c.Closed() || c.handleError(err) {
				// 连接断开或致命错误中断写消息
				break
			}
		} else {
			statistics.SMsgOutput()
		}
	}
}

// handleError 处理上下行消息过程中的错误, 如果是致命错误, 则返回 true
func (c *Client) handleError(err error) bool {
	statistics.SError(err)
	logger.E("handle message error", err.Error())
	if atomic.LoadInt64(&c.id) > 0 {
		Manager.ClientLogout(atomic.LoadInt64(&c.id), c.device)
	}
	return true
}

// Exit 退出客户端
func (c *Client) Exit(code int64) {
	if c.closed.Get() {
		return
	}
	atomic.StoreInt64(&c.id, 0)
	c.closed.Set(true)

	close(c.messages)
	_ = c.conn.Close()
	statistics.SConnExit()
}

// getNextSeq 获取下一个下行消息序列号 sequence
func (c *Client) getNextSeq() int64 {
	seq := c.seq.Get()
	c.seq.Set(seq + 1)
	// TODO SEQ 序列号使用号段模式, 定时持久化号段使用情况
	return seq
}

func (c *Client) Run() {
	logger.I(">>>> client %d running", c.id)
	go c.readMessage()
	go c.writeMessage()
}
