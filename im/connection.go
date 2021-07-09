package im

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

type Connection interface {
	Write(message *Message) error
	Read() (*Message, error)
	Close() error
}

type WsConnection struct {
	options *WsServerOptions
	conn    *websocket.Conn
}

func NewWsConnection(conn *websocket.Conn, options *WsServerOptions) *WsConnection {
	c := new(WsConnection)
	c.conn = conn
	c.options = options
	return c
}

func (c *WsConnection) Write(message *Message) error {
	deadLine := time.Now().Add(c.options.WriteDeadLine)
	_ = c.conn.SetWriteDeadline(deadLine)

	data, err := message.Serialize()
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(1, data)
}

func (c *WsConnection) Read() (*Message, error) {

	_, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return DeserializeMessage(bytes)
}

func (c *WsConnection) error(err error) {
	e := err.Error()
	if strings.HasSuffix(e, "use of closed network conn") {
		return
	}
	_ = c.Close()
	//
}

func (c *WsConnection) Close() error {
	return c.conn.Close()
}
