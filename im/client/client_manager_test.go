package client

import (
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
	return message.NewMessage(0, message.ActionHeartbeat, "").Serialize()
}

func (f *fakeConn) Close() error {
	return nil
}

func TestDefaultManager_ClientConnected(t *testing.T) {
	db.Init()
	dao.Init()
	Manager = NewClientManager()
	c := Manager.ClientConnected(&fakeConn{})
	t.Log(c)
	time.Sleep(time.Second * 20)
}
