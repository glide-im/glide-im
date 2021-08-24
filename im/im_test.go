package im

import (
	"fmt"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/message"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var users []*MockUserConn
var messageCount comm.AtomicInt64 = 0
var mockUid int64 = 1

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
	messageCount.Set(messageCount.Get() + 1)
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

var wg = new(sync.WaitGroup)

//go test -v -run=TestIm -memprofile=mem.out
func TestIm(t *testing.T) {
	userOnline(10_0000)
	sendMessage(5_0000, 100, 100)

	//sec := 0
	//fmt.Println("sec", "\t", "messages")
	//for {
	//	time.Sleep(time.Second)
	//	sec++
	//	fmt.Println(sec, "\t", messageCount)
	//}

	wg.Wait()
	fmt.Println(messageCount)
}

func sendMessage(sender int, count int, duration int64) {
	fmt.Println("start send message")
	max := len(users)
	msg := "hello, how are you, i'm fine, thank you, and you? i'm fine too."
	for i := 0; i < sender; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < count; j++ {
				s := rand.Int63n(duration)
				time.Sleep(time.Millisecond * time.Duration(s))
				from := rand.Intn(max)
				to := rand.Int63n(mockUid + 1)
				go users[from].Send(to, msg)
			}
			wg.Done()
		}()
	}
}

func userOnline(count int) {

	users = make([]*MockUserConn, count, count)
	fmt.Println("load user, count: ", count)
	for i := 0; i < count; i++ {
		con := NewMockUserConn()
		con.uid = client.Manager.ClientConnected(con)
		// user sign in
		client.Manager.ClientSignIn(con.uid, mockUid, -1)
		con.uid = mockUid
		//user[mockUid] = con
		users[i] = con
		mockUid++
	}
	fmt.Println("complete")

}
