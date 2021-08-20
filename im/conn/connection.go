package conn

import (
	"errors"
)

type ConnectionHandler func(conn Connection)

var (
	ErrForciblyClosed   = errors.New("connection was forcibly closed")
	ErrClosed           = errors.New("closed")
	ErrConnectionClosed = errors.New("connection closed")
	ErrBadPackage       = errors.New("bad package data")
)

// Connection expression a network keep-alive connection, WebSocket, tcp etc
type Connection interface {
	// Write write message to the connection.
	Write(message Serializable) error
	// Read read message from the connection.
	Read(message Serializable) error
	// Close close the connection.
	Close() error
}

// ConnectionProxy expression a binder of Connection.
type ConnectionProxy struct {
	conn Connection
}

func (c ConnectionProxy) Write(message Serializable) error {
	return c.conn.Write(message)
}

func (c ConnectionProxy) Read(message Serializable) error {
	return c.conn.Read(message)
}

func (c ConnectionProxy) Close() error {
	return c.conn.Close()
}
