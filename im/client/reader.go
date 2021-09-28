package client

import (
	"go_im/im/conn"
	"go_im/im/message"
)

var messageReader MessageReader = &defaultReader{}

// MessageReader 表示一个从连接中(Connection)读取消息的读取者, 可以用于定义如何从连接中读取并解析消息.
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
