package msgdao

import (
	"go_im/pkg/db"
	"testing"
)

func TestCacheDao_GetIncrUserMsgSeq(t *testing.T) {

	db.Init()
	db.DB.CreateTable(&ChatMessage{})
	db.DB.CreateTable(&OfflineMessage{})
	db.DB.CreateTable(&GroupMessage{})
	db.DB.CreateTable(&GroupMemberMsgState{})
	db.DB.CreateTable(&GroupMessageState{})
	db.DB.CreateTable(&GroupMsgSeq{})
}
