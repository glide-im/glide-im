package client

import (
	"go_im/im/conn"
	"go_im/im/message"
)

var messageReader MessageReader = &defaultReader{}

type MessageReader interface {
	Read(conn conn.Connection) (*message.Message, error)
}

func SetMessageReader(s MessageReader) {
	messageReader = s
}

type defaultReader struct{}

func (d *defaultReader) Read(conn conn.Connection) (*message.Message, error) {
	m := message.Message{}
	err := conn.Read(&m)
	return &m, err
}
