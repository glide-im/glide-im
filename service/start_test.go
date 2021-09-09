package main

import (
	"github.com/gorilla/websocket"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/message"
	"math/rand"
	"testing"
	"time"
)

type MockUserConn struct {
	maxWriteTimeCost int64
	uid              int64
	ms               chan client.SenderChatMessage
}

func NewMockUserConn() *MockUserConn {
	return &MockUserConn{
		maxWriteTimeCost: 10,
		ms:               make(chan client.SenderChatMessage, 100),
	}
}

func (m *MockUserConn) Send(target int64, msg string) {
	m.ms <- client.SenderChatMessage{
		Cid:         1,
		UcId:        1,
		TargetId:    target,
		MessageType: 1,
		Message:     msg,
		SendAt:      dao.Timestamp{},
	}
}

func (m *MockUserConn) Write(s conn.Serializable) error {
	delay := rand.Int63n(m.maxWriteTimeCost)
	time.Sleep(time.Millisecond * time.Duration(delay))
	return nil
}

func (m *MockUserConn) Read(s conn.Serializable) error {
	msg := <-m.ms
	s = message.NewMessage(1, message.ActionChatMessage, msg)
	return nil
}

func (m MockUserConn) Close() error {
	panic("closed")
	return nil
}

func TestClientServer(t *testing.T) {
	runClientService(TypeClientService)
}

func TestClientClient(t *testing.T) {

	go runClientService(TypeApiService)
	time.Sleep(time.Second * 2)
}

func TestClientOnline(t *testing.T) {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	con, _, err := dialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			messageType, p, e := con.ReadMessage()
			if e != nil {
				t.Log(e)
				break
			}
			t.Log(messageType, string(p))
		}

	}()

	_ = con.WriteJSON(message.NewMessage(1, "api.app.echo", ""))
	time.Sleep(time.Hour)
}
