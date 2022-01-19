package conn

import (
	"github.com/gorilla/websocket"
	"net"
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

func (c *WsConnection) Write(data []byte) error {
	deadLine := time.Now().Add(c.options.WriteTimeout)
	_ = c.conn.SetWriteDeadline(deadLine)

	err := c.conn.WriteMessage(websocket.TextMessage, data)
	return c.wrapError(err)
}

func (c *WsConnection) Read() ([]byte, error) {

	deadLine := time.Now().Add(c.options.ReadTimeout)
	_ = c.conn.SetReadDeadline(deadLine)

	msgType, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, c.wrapError(err)
	}

	switch msgType {
	case websocket.TextMessage:
	case websocket.PingMessage:
	case websocket.BinaryMessage:
	default:
		return nil, ErrBadPackage
	}

	return bytes, err
}

func (c *WsConnection) Close() error {
	return c.wrapError(c.conn.Close())
}

func (c *WsConnection) GetConnInfo() *ConnectionInfo {
	c.conn.UnderlyingConn()
	remoteAddr := c.conn.RemoteAddr().(*net.TCPAddr)
	info := ConnectionInfo{
		Ip:   remoteAddr.IP.String(),
		Port: remoteAddr.Port,
		Addr: c.conn.RemoteAddr().String(),
	}
	return &info
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
		return ErrClosed
	}
	if strings.Contains(err.Error(), "use of closed network conn") {
		_ = c.conn.Close()
		return ErrClosed
	}
	if strings.Contains(err.Error(), "i/o timeout") {
		return ErrReadTimeout
	}
	return err
}
