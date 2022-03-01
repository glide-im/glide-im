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

func initDepMock() {
	enqueueMessage = func(uid int64, device int64, message *message.Message) {
		logger.D("%d, %d, %s", uid, device, message.Data)
		atomic.AddInt32(&dispatchCount, 1)
	}
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
	e := UpdateMember(gid, um)
	if e != nil {
		panic(e)
	}
}

func initMock() {

	initDepMock()
	SetInterfaceImpl(NewDefaultManager())
	_ = UpdateGroup(1, Update{Flag: FlagGroupCreate})
}

func TestDefaultManager_dispatch(t *testing.T) {

	initMock()

	var uid []int64
	for i := 0; i < 4; i++ {
		uid = append(uid, int64(i+1))
	}

	initUserMock(1, uid...)

	msg := message.NewChatMessage(1, 1, 1, 1, 1, "HelloWorld", time.Now().Unix())

	for i := 0; i < 4; i++ {
		time.Sleep(time.Millisecond)
		e := dispatch(1, &msg)
		if e != nil {
			t.Error(e)
		}
	}

	t.Log("Msg Dispatch Count:", dispatchCount)
	time.Sleep(time.Second * 3)
}

func TestDefaultManager_dispatch2(t *testing.T) {

	db.Init()
	dao.Init()
	initMock()
	initUserMock(1, 1, 2, 3, 4)

	msg1 := message.NewChatMessage(1, 1, 1, 1, 1, "", time.Now().Unix())
	msg := &msg1
	msg.Mid = 2
	err := dispatch(1, msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	msg.Mid = 3
	err = dispatch(1, msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	msg.Mid = 4
	err = dispatch(1, msg)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 3)
}

func dispatch(gid int64, chatMessage *message.ChatMessage) error {
	return DispatchMessage(gid, chatMessage)
}
