package client

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type Conn struct {
	conn conn.Connection

	closed *comm.AtomicBool
	// buffered channel 40
	messages chan *message.Message
}

func (c *Conn) Close() {
	c.closed.Set(true)
	close(c.messages)
	_ = c.conn.Close()
}

func (c *Conn) Closed() bool {
	return c.closed.Get()
}

func (c *Conn) EnqueueMessage(message *message.Message) {
	select {
	case c.messages <- message:
	default:
		logger.E("Conn.EnqueueMessage", "message chan is full")
	}
}

func (c *Conn) readMessage() (*message.Message, error) {
	msg, err := messageReader.Read(c.conn)
	if err != nil {
		if !c.handleError(err) {
			return nil, err
		}
	}
	return msg, nil
}

func (c *Conn) writeMessage() {
	for msg := range c.messages {
		err := c.conn.Write(msg)
		if err != nil {
			if c.handleError(err) {
				break
			}
		}
	}
}

// handleError return whether fatal error
func (c *Conn) handleError(err error) bool {
	fatalErr := map[error]int{
		conn.ErrForciblyClosed:   0,
		conn.ErrClosed:           0,
		conn.ErrConnectionClosed: 0,
		conn.ErrReadTimeout:      0,
	}
	_, ok := fatalErr[err]
	if ok {
		return true
	}
	return false
}

func (c *Conn) run() {
	go c.writeMessage()
}
