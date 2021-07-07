package im

import (
	"container/list"
)

type Connection struct {
	*mutex

	conn interface{}

	messages *list.List
}

func (c *Connection) Client() {
	defer c.LockFunc()()

}
