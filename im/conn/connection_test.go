package conn

import (
	"github.com/glide-im/glideim/im/message"
	"testing"
	"time"
)

type MockIdleConn struct{}

func (m MockIdleConn) Write(message *message.Message) error { return nil }
func (m MockIdleConn) Read() (*message.Message, error) {
	time.Sleep(time.Hour)
	return nil, nil
}
func (m MockIdleConn) Close() error { return nil }

type Tt interface {
	Fn(t string)
}

type Impl struct {
	s string
}

func (i *Impl) Fn(to string) {
	i.s = to
}

func tf(i Tt) {
	i.Fn("123")
	i = &Impl{s: "2"}
}

func TestNewWsConnection(t *testing.T) {
	i := &Impl{s: "1"}
	tf(i)
	t.Log(i)
}
