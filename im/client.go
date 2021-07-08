package im

import "time"

// Client represent a user client connection
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

func (c *Client) ReadMessage() {
	for {
		message, err := c.conn.Read()
		if err != nil {
			//
			continue
		}
		if message.action.IsApi() {

		}
	}
}

func (c *Client) WriteMessage() {
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
	go c.ReadMessage()
	go c.WriteMessage()
}
