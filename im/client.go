package im

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/entity"
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
	messages chan *entity.Message

	seq *comm.AtomicInt64

	heartbeat *time.Timer
}

func NewClient(conn conn.Connection, connUid int64) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = comm.NewAtomicBool(false)
	client.messages = make(chan *entity.Message, 40)
	client.time = time.Now()
	client.uid = connUid
	client.seq = new(comm.AtomicInt64)
	// TODO 优化内存
	client.heartbeat = time.AfterFunc(HeartbeatDuration, client.onDeath)
	client.heartbeat.Stop()
	return client
}

func (c *Client) EnqueueMessage(message *entity.Message) {
	comm.Slog.I("EnqueueMessage(uid=%d, %s): %v", c.uid, message.Action, message)
	if c.closed.Get() {
		comm.Slog.W("connection closed, cannot enqueue message")
		return
	}
	if message.Seq <= 0 {
		message.Seq = c.getNextSeq()
	}
	select {
	case c.messages <- message:
		break
	default:
		comm.Slog.E("Client.EnqueueMessage", "message chan is full")
	}
}

func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			comm.Slog.D("Recover: conn read message error: %v", err)
		}
	}()

	comm.Slog.I("start read message")
	for {
		message := &entity.Message{}
		err := c.conn.Read(message)
		if err != nil {
			if !c.handleError(-1, err) {
				continue
			}
			break
		}
		if message.Action.Contains(entity.ActionApi) {
			ApiManager.Handle(c.uid, message)
		} else if message.Action.Contains(entity.ActionMessage) {
			err = c.dispatch(message)
		} else if message.Action == entity.ActionHeartbeat {
			c.heartbeat.Reset(HeartbeatDuration)
		} else {
			// unknown action
		}
		if err != nil {
			if !c.handleError(message.Seq, err) {
				continue
			}
			break
		}
	}
}

func (c *Client) writeMessage() {
	comm.Slog.I("start write message")

	for message := range c.messages {
		err := c.conn.Write(message)
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
	}
	_, ok := fatalErr[err]
	if ok {
		ClientManager.UserLogout(c.uid)
		return true
	}
	c.messages <- entity.NewMessage(seq, entity.ActionNotify, err.Error())
	return false
}

func (c *Client) onDeath() {
	// TODO
}

func (c *Client) dispatch(message *entity.Message) error {
	switch message.Action {
	case entity.ActionChatMessage:
		return ClientManager.DispatchMessage(c.uid, message)
	case entity.ActionGroupMessage:
		return GroupManager.DispatchMessage(c.uid, message)
	default:
		// unknown message type
	}
	return nil
}

func (c *Client) Exit() {
	c.closed.Set(true)
	close(c.messages)
	_ = c.conn.Close()
}

func (c *Client) getNextSeq() int64 {
	seq := c.seq.Get()
	c.seq.Set(seq + 1)
	return seq
}

func (c *Client) Run() {
	comm.Slog.I("///////////////////////// connection running /////////////////////////////")
	go c.readMessage()
	go c.writeMessage()
	c.heartbeat.Reset(HeartbeatDuration)
}
