package group

import (
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"sync/atomic"
	"testing"
	"time"
)

var dispatchCount int32 = 0

type MockClientManager struct {
}

func (m *MockClientManager) EnqueueMessage(uid int64, device int64, message *message.Message) {
	logger.D("%d, %d, %s", uid, device, message.Data)
	atomic.AddInt32(&dispatchCount, 1)
}

func initDepMock() {
	EnqueueMessage = &MockClientManager{}
}

func initUserMock(gid int64, uid ...int64) {
	var um []MemberUpdate
	for _, i := range uid {
		um = append(um, MemberUpdate{
			Uid:  i,
			Flag: FlagMemberAdd,
		})

		um = append(um, MemberUpdate{
			Uid:  i,
			Flag: FlagMemberOnline,
		})
	}
	e := Manager.UpdateMember(gid, um)
	if e != nil {
		panic(e)
	}
}

func initMock() {

	initDepMock()

	Manager = NewDefaultManager()

	_ = Manager.UpdateGroup(1, Update{Flag: FlagGroupCreate})
}

func TestDefaultManager_DispatchMessage(t *testing.T) {

	initMock()

	var uid []int64
	for i := 0; i < 4; i++ {
		uid = append(uid, int64(i+1))
	}

	initUserMock(1, uid...)

	msg := &message.UpChatMessage{
		Mid:     1,
		CSeq:    1,
		From:    1,
		To:      1,
		Type:    1,
		Content: "HelloWorld",
		CTime:   time.Now().Unix(),
	}

	for i := 0; i < 4; i++ {
		time.Sleep(time.Millisecond)
		e := Manager.DispatchMessage(1, msg)
		if e != nil {
			t.Error(e)
		}
	}

	t.Log("Msg Dispatch Count:", dispatchCount)
	time.Sleep(time.Second * 3)
}

func TestDefaultManager_DispatchMessage2(t *testing.T) {

	db.Init()
	dao.Init()
	initMock()
	initUserMock(1, 1, 2, 3, 4)

	msg := &message.UpChatMessage{
		Mid:     1,
		CSeq:    1,
		From:    1,
		To:      1,
		Type:    1,
		Content: "HelloWorld",
		CTime:   time.Now().Unix(),
	}
	msg.Mid = 2
	err := Manager.DispatchMessage(1, msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	msg.Mid = 3
	err = Manager.DispatchMessage(1, msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	msg.Mid = 4
	err = Manager.DispatchMessage(1, msg)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 3)
}
