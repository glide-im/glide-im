package im

import (
	"fmt"
	"go_im/im/entity"
	"time"
)

// Client represent a user client conn
type Client struct {
	conn Connection

	uid      int64
	deviceId int64
	time     time.Time

	messages chan *entity.Message
	closed   []chan interface{}
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
	return false
}

func (c *Client) Offline(reason string) {

}

func (c *Client) readMessage() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for {
		message, err := c.conn.Read()
		if err != nil {
			c.handleError(-1, err)
			continue
		}
		if message.Action&entity.MaskActionApi != 0 {
			err = Api.Handle(c, message)
		} else if message.Action&entity.MaskActionMessage != 0 {
			c.handleMessage(message)
		} else if message.Action == entity.ActionHeartbeat {
			c.handleHeartbeat(message)
		}
		if err != nil {
			c.handleError(message.Seq, err)
		}
	}
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

func (c *Client) handleError(seq int64, err error) {
	c.messages <- &entity.Message{
		Seq:  seq,
		Data: []byte(err.Error()),
	}
}

func (c *Client) handleHeartbeat(message *entity.Message) {

}

func (c *Client) handleMessage(message *entity.Message) {

}

func (c *Client) Run() {
	go c.readMessage()
	go c.writeMessage()
}
