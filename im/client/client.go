package client

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/pkg/logger"
	"time"
)

const HeartbeatDuration = time.Second * 30

// Client represent a user conn conn
type Client struct {
	conn conn.Connection

	uid      int64
	deviceId int64
	time     time.Time
	closed   *comm.AtomicBool

	// buffered channel 40
	messages chan *message.Message

	seq *comm.AtomicInt64

	heartbeat *time.Timer
}

func NewClient(conn conn.Connection, connUid int64) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = comm.NewAtomicBool(false)
	client.messages = make(chan *message.Message, 40)
	client.time = time.Now()
	client.uid = connUid
	client.seq = new(comm.AtomicInt64)
	// TODO 优化内存
	client.heartbeat = time.AfterFunc(HeartbeatDuration, client.onDeath)
	client.heartbeat.Stop()
	return client
}

func (c *Client) SignIn(uid int64, device int64) {
	c.uid = uid
	c.deviceId = device
}

func (c *Client) Id() int64 {
	return c.uid
}

func (c *Client) Closed() bool {
	return c.closed.Get()
}

func (c *Client) EnqueueMessage(message *message.Message) {
	logger.I("EnqueueMessage(uid=%d, %s): %v", c.uid, message.Action, message)
	if c.closed.Get() {
		logger.W("connection closed, cannot enqueue message")
		return
	}
	if message.Seq <= 0 {
		message.Seq = c.getNextSeq()
	}
	select {
	case c.messages <- message:
	default:
		logger.E("Client.EnqueueMessage", "message chan is full")
	}
}

func (c *Client) readMessage() {
	defer func() {
		//err := recover()
		//if err != nil {
		//	comm.D("Recover: conn read message error: %v", err)
		//}
	}()

	logger.I("start read message")
	for {
		msg, err := messageReader.Read(c.conn)
		if err != nil {
			if !c.handleError(-1, err) {
				continue
			}
			break
		}
		if msg.Action.Contains(message.ActionApi) {
			Manager.Api(c.uid, msg)
		} else if msg.Action.Contains(message.ActionMessage) {
			err = Manager.DispatchMessage(c.uid, msg)
		} else if msg.Action == message.ActionHeartbeat {
			c.heartbeat.Reset(HeartbeatDuration)
		} else {
			// unknown action
		}
		if err != nil {
			if !c.handleError(msg.Seq, err) {
				continue
			}
			break
		}
	}
}

func (c *Client) writeMessage() {
	logger.I("start write message")

	for msg := range c.messages {
		err := c.conn.Write(msg)
		if err != nil {
			if c.handleError(-1, err) {
				break
			}
		}
	}
}

// handleError return whether fatal error
func (c *Client) handleError(seq int64, err error) bool {

	fatalErr := map[error]int{
		conn.ErrForciblyClosed:   0,
		conn.ErrClosed:           0,
		conn.ErrConnectionClosed: 0,
		conn.ErrReadTimeout:      0,
	}
	_, ok := fatalErr[err]
	if ok {
		Manager.UserLogout(c.uid)
		return true
	}
	logger.E("err", err.Error())
	c.EnqueueMessage(message.NewMessage(seq, "notify", err.Error()))
	return false
}

func (c *Client) onDeath() {
	// TODO
}

func (c *Client) Exit() {
	if c.closed.Get() {
		return
	}
	c.closed.Set(true)
	close(c.messages)
	c.heartbeat.Stop()
	_ = c.conn.Close()
}

func (c *Client) getNextSeq() int64 {
	seq := c.seq.Get()
	c.seq.Set(seq + 1)
	return seq
}

func (c *Client) Run() {
	logger.I("///////////////////////// connection running /////////////////////////////")
	go c.readMessage()
	go c.writeMessage()
	c.heartbeat.Reset(HeartbeatDuration)
}
