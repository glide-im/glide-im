package im

import (
	"go_im/im/entity"
	"time"
)

// Client represent a user client conn
type Client struct {
	conn Connection

	uid      int64
	deviceId int64
	time     time.Time
	closed   *AtomicBool

	groups   []*Group
	messages chan *entity.Message
}

func NewClient(conn Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.closed = NewAtomicBool(false)
	client.messages = make(chan *entity.Message, 200)
	client.time = time.Now()

	return client
}

func (c *Client) AddGroup(group *Group) {
	c.groups = append(c.groups, group)
}

// EnqueueMessage enqueue blocking message channel
func (c *Client) EnqueueMessage(message *entity.Message) {
	logger.I("enqueue new message %v", message)
	if c.closed.Get() {
		logger.W("connection closed, cannot enqueue message")
		return
	}
	c.messages <- message
}

func (c *Client) SignIn(uid int64, deviceId int64) {
	logger.I("client sign in")
	c.uid = uid
	c.deviceId = deviceId
	ClientManager.ClientSignIn(c)
	logger.D("client sign in uid=%d", uid)
}

func (c *Client) SignOut(reason string) {
	logger.I("client sign out")
	c.uid = 0
	for _, group := range c.groups {
		group.Unsubscribe(c.uid)
	}
	ClientManager.ClientSignOut(c)
	_ = c.conn.Close()
	logger.D("connection closed uid=%d, reason=%d", c.uid, reason)
}

func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.D("client read message error: %v", err)
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
		if message.Action&entity.MaskActionApi != 0 {
			err = Api.Handle(c, message)
		} else if message.Action&entity.MaskActionMessage != 0 {
			err = c.handleMessage(message)
		} else if message.Action == entity.ActionHeartbeat {
			c.handleHeartbeat(message)
		}
		if err != nil {
			if !c.handleError(message.Seq, err) {
				continue
			}
			break
		}
	}
	c.uid = 0
}

func (c *Client) writeMessage() {
	logger.I("start write message")
	for {
		select {
		// blocking write
		case message := <-c.messages:
			err := c.conn.Write(message)
			if err != nil {
				logger.E("client write message error", err)
			}
		}
	}
}

// handleError return whether fatal error
func (c *Client) handleError(seq int64, err error) bool {

	if err == ErrForciblyClosed {
		c.closed.Set(true)
		c.SignOut("forcibly closed")
		logger.D("uid=%d forcibly closed", c.uid)
		return true
	}

	//c.messages <- entity.NewErrMessage(seq, err)

	return false
}

func (c *Client) handleHeartbeat(message *entity.Message) {

}

func (c *Client) handleMessage(message *entity.Message) error {
	switch message.Action {
	case entity.ActionChatMessage:
		return ClientManager.DispatchMessage(c.uid, message)
	case entity.ActionGroupMessage:
		return GroupManager.DispatchMessage(c, message)
	}
	return nil
}

func (c *Client) Run() {
	logger.D("///////////////////////// connection running /////////////////////////////")
	go c.readMessage()
	go c.writeMessage()
}
