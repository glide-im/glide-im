package im

import (
	"go_im/im/entity"
	"strings"
	"time"
)

// Client represent a user client conn
type Client struct {
	conn Connection

	uid      int64
	deviceId int64
	time     time.Time
	closed   *AtomicBool

	messages chan *entity.Message

	seq *AtomicInt64
}

func NewClient(conn Connection, connUid int64) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = NewAtomicBool(false)
	client.messages = make(chan *entity.Message, 200)
	client.time = time.Now()
	client.uid = connUid
	client.seq = new(AtomicInt64)
	return client
}

// EnqueueMessage enqueue blocking message channel
func (c *Client) EnqueueMessage(message *entity.Message) {
	logger.I("EnqueueMessage(uid=%d, %s): %v", c.uid, message.Action, message)
	if c.closed.Get() {
		logger.W("connection closed, cannot enqueue message")
		return
	}
	if message.Seq <= 0 {
		message.Seq = c.getNextSeq()
	}
	c.messages <- message
}

func (c *Client) SignOut(reason string) {
	if c.closed.Get() {
		logger.E("Client.SignOut", "client has already sign out")
		return
	}
	logger.I("client sign out uid=%d, reason=%s", c.uid, reason)
	ClientManager.UserLogout(c.uid)
	c.Exit()
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

func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("Recover: client read message error: %v", err)
		}
	}()

	logger.I("start read message")
	for {
		message, err := c.conn.Read()
		if err != nil {
			if !c.handleError(-1, err) {
				continue
			}
			break
		}
		logger.D("NewMessage(uid=%d, %s): %s", c.uid, message.Action, message)
		if message.Action&entity.MaskActionApi != 0 {
			err = Api.Handle(c.uid, message)
		} else if message.Action&entity.MaskActionMessage != 0 {
			err = c.handleMessage(message)
		} else if message.Action == entity.ActionHeartbeat {
			c.handleHeartbeat(message)
		} else {
			// echo
			m, _ := message.Serialize()
			c.EnqueueMessage(entity.NewSimpleMessage(1, entity.ActionEcho, string(m)))
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
	logger.I("start write message")

	hello := entity.NewMessage2(time.Now().Unix(), entity.ActionAck, "hello")
	c.EnqueueMessage(hello)

	for {
		select {
		// blocking write
		case message := <-c.messages:
			err := c.conn.Write(message)
			if err != nil {
				logger.E("client write message error", err)
				c.handleError(-1, err)
				if c.closed.Get() {
					break
				}
			}
		}
		if c.closed.Get() {
			logger.D("write message break len=%d", len(c.messages))
			break
		}
	}
}

// handleError return whether fatal error
func (c *Client) handleError(seq int64, err error) bool {

	if strings.Contains(err.Error(), "use of closed network connection") {
		c.SignOut("connection closed")
		return true
	}

	if err == ErrForciblyClosed || err == ErrClosed {
		c.SignOut("client forcibly closed")
		return true
	}

	c.messages <- entity.NewErrMessage(seq, err)

	return false
}

func (c *Client) handleHeartbeat(message *entity.Message) {

}

func (c *Client) handleMessage(message *entity.Message) error {
	switch message.Action {
	case entity.ActionChatMessage:
		return ClientManager.DispatchMessage(c.uid, message)
	case entity.ActionGroupMessage:
		return GroupManager.DispatchMessage(c.uid, message)
	}
	return nil
}

func (c *Client) Run() {
	logger.D("///////////////////////// connection running /////////////////////////////")
	go c.readMessage()
	go c.writeMessage()
}
