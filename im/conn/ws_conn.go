package conn

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

type WsConnection struct {
	options *WsServerOptions
	conn    *websocket.Conn
}

func NewWsConnection(conn *websocket.Conn, options *WsServerOptions) *WsConnection {
	c := new(WsConnection)
	c.conn = conn
	c.options = options
	c.conn.SetCloseHandler(func(code int, text string) error {
		return ErrClosed
	})
	return c
}

func (c *WsConnection) Write(message Serializable) error {
	deadLine := time.Now().Add(c.options.WriteTimeout)
	_ = c.conn.SetWriteDeadline(deadLine)

	data, err := message.Serialize()
	if err != nil {
		return c.wrapError(err)
	}
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	return c.wrapError(err)
}

func (c *WsConnection) Read(message Serializable) error {

	deadLine := time.Now().Add(c.options.ReadTimeout)
	_ = c.conn.SetReadDeadline(deadLine)

	msgType, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return c.wrapError(err)
	}

	switch msgType {
	case websocket.TextMessage:
	case websocket.PingMessage:
	case websocket.BinaryMessage:
	default:
		return ErrBadPackage
	}

	return message.Deserialize(bytes)
}

func (c *WsConnection) Close() error {
	return c.wrapError(c.conn.Close())
}

func (c *WsConnection) wrapError(err error) error {
	if err == nil {
		return nil
	}
	if websocket.IsUnexpectedCloseError(err) {
		return ErrClosed
	}
	if websocket.IsCloseError(err) {
		return ErrClosed
	}
	if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
		_ = c.conn.Close()
		return ErrForciblyClosed
	}
	if strings.Contains(err.Error(), "use of closed network conn") {
		_ = c.conn.Close()
		return ErrConnectionClosed
	}
	if strings.Contains(err.Error(), "i/o timeout") {
		return ErrReadTimeout
	}
	return err
}
