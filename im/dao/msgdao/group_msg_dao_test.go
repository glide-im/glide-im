package msgdao

import (
	"go_im/pkg/db"
	"testing"
	"time"
)

var dao GroupMsgDao = groupMsgDaoImpl{}

func init() {
	db.Init()
}

func TestGroupMsgDao_CreateGroupMsgSeq(t *testing.T) {
	err := CreateGroupMsgSeq(1, 100)
	if err != nil {
		t.Error(err)
	}
}

func TestGroupMsgDao_GetGroupMsgSeq(t *testing.T) {
	seq, err := GetGroupMsgSeq(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(seq)
}

func TestGroupMsgDao_GetGroupMessageSeqAfter(t *testing.T) {
	ms, err := GetGroupMessageSeqAfter(1, 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(ms)
}

func TestGroupMsgDao_AddGroupMessage(t *testing.T) {
	err := AddGroupMessage(&GroupMessage{
		MID:     1231241238,
		Seq:     28,
		To:      1,
		From:    2,
		Type:    1,
		SendAt:  time.Now().Unix(),
		Content: "hello world",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGroupMsgDao_GetGroupMessage(t *testing.T) {
	message, err := GetGroupMessage(1231241235)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func TestGroupMsgDao_GetGroupMessageState(t *testing.T) {
	state, err := dao.GetGroupMessageState(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(state)
}

func TestGroupMsgDao_UpdateGroupMessageState(t *testing.T) {
	err := dao.UpdateGroupMessageState(1, 1, time.Now().Unix(), 2)
	if err != nil {
		t.Error(err)
	}
}

func TestGroupMsgDao_CreateGroupMemberMsgState(t *testing.T) {
	err := dao.CreateGroupMemberMsgState(2, 1)
	if err != nil {
		t.Log(err)
	}
}

func TestGroupMsgDao_GetGroupMemberMsgState(t *testing.T) {
	state, err := dao.GetGroupMemberMsgState(2, 1)
	if err != nil {
		t.Log(err)
	}
	t.Log(state)
}

func TestGroupMsgDao_UpdateGroupMemberMsgState(t *testing.T) {
	err := dao.UpdateGroupMemberMsgState(2, 1, 1, 1)
	if err != nil {
		t.Log(err)
	}
}
