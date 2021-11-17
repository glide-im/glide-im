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

// MessageHandleFunc 所有客户端消息都传递到该函数处理
var MessageHandleFunc func(from int64, device int64, message *message.Message) = nil

var tw = timingwheel.NewTimingWheel(time.Millisecond*500, 3, 20)

// HeartbeatDuration 心跳间隔, 默认5分钟
const HeartbeatDuration = time.Minute * 5

// IClient 表示一个客户端, 用于管理连接状态, 连接 id, 消息收发
type IClient interface {

	// SetID 设置该客户端标识 ID
	SetID(id int64, device int64)

	// Closed 返回该客户端连接是否已关闭
	Closed() bool

	// EnqueueMessage 将消息放入到客户端消息队列中
	EnqueueMessage(message *message.Message)

	// Exit 退出客户端, 关闭连接等
	Exit()

	// Run 开始收发消息客户端连接的消息
	Run()
}

// Client represent a user conn conn
type Client struct {
	conn conn.Connection

	// id 唯一标识
	id int64
	// device 设备标识
	device    int64
	connectAt time.Time
	// closed 连接是否关闭
	closed *comm.AtomicBool

	// messages 带缓冲的下行消息管道, 缓冲大小40
	messages chan *message.Message
	// readClose 关闭或写入则停止读
	readClose chan struct{}
	// writeClose 关闭或写入则停止写
	writeClose chan struct{}

	// hb 心跳倒计时
	hb *timingwheel.Task

	// seq 服务器下行消息递增序列号
	seq *comm.AtomicInt64
}

func newClient(conn conn.Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = comm.NewAtomicBool(false)
	// 大小为 40 的缓冲管道, 防止短时间消息过多如果网络连接 output 不及时会造成程序阻塞, 可以适当调整
	client.messages = make(chan *message.Message, 40)
	client.connectAt = time.Now()
	client.readClose = make(chan struct{})
	client.writeClose = make(chan struct{})
	client.seq = new(comm.AtomicInt64)
	client.hb = tw.After(HeartbeatDuration)
	return client
}

// SetID 设置 id 标识及设备标识
func (c *Client) SetID(id int64, device int64) {
	atomic.StoreInt64(&c.id, id)
	atomic.StoreInt64(&c.device, device)
}

func (c *Client) Closed() bool {
	return c.closed.Get()
}

// EnqueueMessage 放入下行消息队列
func (c *Client) EnqueueMessage(message *message.Message) {
	if c.Closed() {
		return
	}
	logger.I("EnqueueMessage(id=%d, %s): %v", c.id, message.Action, message)
	if message.Seq < 0 {
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
	readChan, done := messageReader.ReadCh(c.conn)

	defer func() {
		err := recover()
		if err != nil {
			logger.E("read message error", err)
		}
	}()

	for {
		select {
		case <-c.readClose:
			goto STOP
		case <-c.hb.C:
			// TODO 处理心跳超时
			logger.W("heartbeat timout")
		case msg, ok := <-readChan:
			if !ok {
				goto STOP
			}
			if msg.err != nil {
				if c.Closed() || c.handleError(msg.err) {
					// 连接断开或致命错误中断读消息
					goto STOP
				}
				continue
			}
			c.hb.Cancel()
			c.hb = tw.After(HeartbeatDuration)
			// 统一处理消息函数
			MessageHandleFunc(c.id, c.device, msg.m)
			msg.Recycle()
		}
	}
STOP:
	close(done)
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
		case <-c.writeClose:
			goto STOP
		case m := <-c.messages:
			b, err := m.Serialize()
			if err != nil {
				logger.E("serialize output message", err)
				continue
			}
			err = c.conn.Write(b)
			if err != nil {
				if c.Closed() || c.handleError(err) {
					// 连接断开或致命错误中断写消息
					goto STOP
				}
			} else {
				statistics.SMsgOutput()
			}
			//case <-timeout:
			// TODO write message time is too long, slow client
		}
	}
STOP:
}

// handleError 处理上下行消息过程中的错误, 如果是致命错误, 则返回 true
func (c *Client) handleError(err error) bool {
	statistics.SError(err)
	logger.E("handle message error: %s", err.Error())
	if atomic.LoadInt64(&c.id) > 0 {
		Manager.ClientLogout(atomic.LoadInt64(&c.id), c.device)
	}
	return true
}

// Exit 退出客户端
func (c *Client) Exit() {
	// TODO 先关闭下行消息队列写入, 真正退出前先将下行队列然后所有消息写完
	if c.closed.Get() {
		return
	}
	c.closed.Set(true)
	atomic.StoreInt64(&c.id, 0)

	close(c.readClose)
	close(c.writeClose)
	//close(c.messages)
	_ = c.conn.Close()
	statistics.SConnExit()
}

// getNextSeq 获取下一个下行消息序列号 sequence
func (c *Client) getNextSeq() int64 {
	seq := c.seq.Get()
	c.seq.Set(seq + 1)
	return seq
}

func (c *Client) Run() {
	logger.I(">>>> client %d running", c.id)
	go c.readMessage()
	go c.writeMessage()
}
