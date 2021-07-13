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
	closed   bool

	messages chan *entity.Message
}

func NewClient(conn Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.messages = make(chan *entity.Message, 200)
	client.time = time.Now()

	return client
}

func (c *Client) EnqueueMessage(message *entity.Message) {
	c.messages <- message
}

func (c *Client) IsOnline() bool {
	return c.uid != 0
}

func (c *Client) Close(reason string) {
	c.uid = 0
	_ = c.conn.Close()
	logger.D("connection closed uid=%d", c.uid)
}

func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			logger.E("client read message error", err.(error))
		}
	}()

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
			c.handleMessage(message)
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
	for {
		select {
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
		c.closed = true
		logger.D("uid=%d forcibly closed", c.uid)
		return true
	}

	c.messages <- entity.NewErrMessage(seq, err)

	return false
}

func (c *Client) handleHeartbeat(message *entity.Message) {

}

func (c *Client) handleMessage(message *entity.Message) {
	switch message.Action {
	case entity.ActionChatMessage:
	case entity.ActionGroupMessage:
		GroupManager.DispatchMessage(c, message)
	}
}

func (c *Client) Run() {
	go c.readMessage()
	go c.writeMessage()
	logger.D("new connection")
}
