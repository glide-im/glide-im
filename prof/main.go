package main

import (
	"fmt"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/message"
	"math/rand"
	"time"
)

var users []*MockUserConn
var messageCount comm.AtomicInt64 = 0
var mockUid int64 = 1

type MockUserConn struct {
	maxWriteTimeCost int64
	uid              int64
	ms               chan message.SenderChatMessage
}

func NewMockUserConn() *MockUserConn {
	return &MockUserConn{
		maxWriteTimeCost: 10,
		ms:               make(chan message.SenderChatMessage, 100),
	}
}

func (m MockUserConn) Send(target int64, msg string) {
	messageCount.Set(messageCount.Get() + 1)
	m.ms <- message.SenderChatMessage{
		Cid:         1,
		UcId:        1,
		TargetId:    target,
		MessageType: 1,
		Message:     msg,
		SendAt:      time.Now().Unix(),
	}
}

func (m MockUserConn) Write(message conn.Serializable) error {
	delay := rand.Int63n(m.maxWriteTimeCost)
	time.Sleep(time.Millisecond * time.Duration(delay))
	return nil
}

func (m MockUserConn) Read(s conn.Serializable) error {
	msg := <-m.ms
	message.NewMessage(1, message.ActionChatMessage, msg)
	return nil
}

func (m MockUserConn) Close() error {
	panic("closed")
	return nil
}

func main() {
	userOnline(1000)
	sendMessage(1000, 50, 50)

	sec := 0
	fmt.Println("sec", "\t", "messages")
	for messageCount.Get() < 40_000 {
		time.Sleep(time.Second)
		sec++
		fmt.Println(sec, "\t", messageCount.Get())
	}
}

func sendMessage(sender int, count int, duration int64) {
	fmt.Println("start send message")
	max := len(users)
	msg := "hello, how are you, i'm fine, thank you, and you? i'm fine too."
	for i := 0; i < sender; i++ {
		go func() {
			for j := 0; j < count; j++ {
				s := rand.Int63n(duration)
				time.Sleep(time.Millisecond * time.Duration(s))
				from := rand.Intn(max)
				to := rand.Int63n(mockUid + 1)
				go users[from].Send(to, msg)
			}
		}()
	}
}

func userOnline(count int) {

	users = make([]*MockUserConn, count, count)
	fmt.Println("load user, count: ", count)
	for i := 0; i < count; i++ {
		conn := NewMockUserConn()
		conn.uid = client.Manager.ClientConnected(conn)
		// user sign in
		client.Manager.ClientSignIn(conn.uid, mockUid, -1)
		conn.uid = mockUid
		//user[mockUid] = conn
		users[i] = conn
		mockUid++
	}
	fmt.Println("complete")

}
