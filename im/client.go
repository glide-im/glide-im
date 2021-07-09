package im

import (
	"fmt"
	"time"
)

// Client represent a user client conn
type Client struct {
	conn Connection

	uid      int64
	deviceId int64
	time     time.Time

	messages chan *Message
}

func NewClient(conn Connection) *Client {
	client := new(Client)
	client.conn = conn
	client.messages = make(chan *Message, 200)
	client.time = time.Now()

	return client
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
			//
			continue
		}
		if message.Action.IsApi() {

		} else if message.Action.IsHeartbeat() {

		} else if message.Action.IsMessage() {

		}
	}
}

func (c *Client) deliver() {

}

func (c *Client) IsOnline() bool {
	return false
}

func (c *Client) writeMessage() {
	for {
		select {
		case message := <-c.messages:
			err := c.conn.Write(message)
			if err != nil {
				//
			}
		}
	}
}

func (c *Client) Run() {
	go c.readMessage()
	go c.writeMessage()
}
