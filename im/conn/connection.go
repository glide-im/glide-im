package conn

import (
	"errors"
)

var (
	ErrForciblyClosed   = errors.New("connection was forcibly closed")
	ErrClosed           = errors.New("closed")
	ErrConnectionClosed = errors.New("connection closed")
	ErrBadPackage       = errors.New("bad package data")
	ErrReadTimeout      = errors.New("i/o timeout")
)

// Connection expression a network keep-alive connection, WebSocket, tcp etc
type Connection interface {
	// Write message to the connection.
	Write(data []byte) error
	// Read message from the connection.
	Read() ([]byte, error)
	// Close the connection.
	Close() error
}

// ConnectionProxy expression a binder of Connection.
type ConnectionProxy struct {
	conn Connection
}

func (c ConnectionProxy) Write(data []byte) error {
	return c.conn.Write(data)
}

func (c ConnectionProxy) Read() ([]byte, error) {
	return c.conn.Read()
}

func (c ConnectionProxy) Close() error {
	return c.conn.Close()
}
