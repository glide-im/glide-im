package im

import (
	"github.com/gorilla/websocket"
	"go_im/im/entity"
	"strings"
	"time"
)

type Connection interface {
	Write(message *entity.Message) error
	Read() (*entity.Message, error)
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

func (c *WsConnection) Write(message *entity.Message) error {
	deadLine := time.Now().Add(c.options.WriteDeadLine)
	_ = c.conn.SetWriteDeadline(deadLine)

	data, err := message.Serialize()
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(1, data)
}

func (c *WsConnection) Read() (*entity.Message, error) {

	_, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return entity.DeserializeMessage(bytes)
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
