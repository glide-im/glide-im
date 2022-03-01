package client

import (
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/db"
	"math/rand"
	"testing"
	"time"
)

type fakeConn struct{}

func (f *fakeConn) Write(data []byte) error {
	return nil
}

func (f *fakeConn) Read() ([]byte, error) {
	time.Sleep(time.Second * time.Duration(3+rand.Int63n(2)))
	newMessage := message.NewMessage(0, message.ActionHeartbeat, "")
	return codec.Encode(newMessage)
}

func (f *fakeConn) Close() error {
	return nil
}

func (f *fakeConn) GetConnInfo() *conn.ConnectionInfo {
	return nil
}

func TestDefaultManager_ClientConnected(t *testing.T) {
	db.Init()
	dao.Init()
	manager := NewDefaultManager()
	manager = manager
	c := manager.ClientConnected(&fakeConn{})
	t.Log(c)
	time.Sleep(time.Second * 20)
}
