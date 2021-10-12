package client

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/pkg/logger"
	"time"
)

// MessageHandler 用于处理客户端消息
type MessageHandler func(from int64, device int64, message *message.Message)

// MessageHandleFunc 所有客户端消息都传递到该函数处理
var MessageHandleFunc MessageHandler = nil

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

	id     int64
	device int64
	time   time.Time
	closed *comm.AtomicBool

	// buffered channel 40
	messages chan *message.Message

	seq *comm.AtomicInt64
}

func NewClient(conn conn.Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = comm.NewAtomicBool(false)
	// 大小为 40 的缓冲管道, 防止短时间消息过多如果网络连接 output 不及时会造成程序阻塞, 可以适当调整
	client.messages = make(chan *message.Message, 40)
	client.time = time.Now()
	client.seq = new(comm.AtomicInt64)
	return client
}

func (c *Client) SetID(id int64, device int64) {
	c.id = id
	c.device = device
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
		// 消息 chan 缓冲溢出, 这条消息将被丢弃
		logger.E("Client.EnqueueMessage", "message chan is full")
	}
}

// readMessage 开始从 Connection 中读取消息
func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("Client read message error: %v", err)
		}
	}()

	for {
		msg, err := messageReader.Read(c.conn)
		if err != nil {
			if c.Closed() || c.handleError(-1, err) {
				// 连接断开或致命错误中断读消息
				break
			}
			continue
		}
		if msg.Action == message.ActionHeartbeat {

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
	}()

	for msg := range c.messages {
		err := c.conn.Write(msg)
		if err != nil {
			if c.Closed() || c.handleError(-1, err) {
				// 连接断开或致命错误中断写消息
				break
			}
		}
	}
}

// handleError 处理上下行消息过程中的错误, 如果是致命错误, 则返回 true
func (c *Client) handleError(seq int64, err error) bool {

	fatalErr := map[error]int{
		conn.ErrForciblyClosed:   0,
		conn.ErrClosed:           0,
		conn.ErrConnectionClosed: 0,
		conn.ErrReadTimeout:      0,
	}
	_, ok := fatalErr[err]
	if ok {
		logger.D("handle message fatal error: %s", err.Error())
		if c.id >= 0 {
			Manager.ClientLogout(c.id, c.device)
		}
		return true
	}
	logger.E("handle message error", err.Error())
	c.EnqueueMessage(message.NewMessage(seq, "notify", err.Error()))
	return false
}

// Exit 退出客户端
func (c *Client) Exit() {
	if c.closed.Get() {
		return
	}
	c.id = 0
	c.closed.Set(true)

	close(c.messages)
	_ = c.conn.Close()
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
