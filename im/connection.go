package im

import (
	"errors"
	"github.com/gorilla/websocket"
	"go_im/im/entity"
	"strings"
	"time"
)

var (
	ErrForciblyClosed = errors.New("connection was forcibly closed")
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
	c.conn.SetCloseHandler(func(code int, text string) error {
		logger.D("closed")
		return nil
	})
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
		if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
			_ = c.conn.Close()
			err = ErrForciblyClosed
		}
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
