package msgdao

import (
	"sync/atomic"
	"time"
)

type chatMsgMock struct {
	s time.Duration
}

func MockChatMsg(queryDelay time.Duration) {
	ChatMsgDaoImpl = &chatMsgMock{s: queryDelay}
}

func (c *chatMsgMock) GetChatMessage(mid ...int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) GetChatMessagesBySession(uid1, uid2 int64, beforeMid int64, pageSize int) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) GetRecentChatMessagesBySession(uid1, uid2 int64, pageSize int) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) GetRecentChatMessages(uid int64, afterTime int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) AddChatMessage(message *ChatMessage) (bool, error) {
	time.Sleep(c.s)
	return true, nil
}

func (c *chatMsgMock) UpdateChatMessageStatus(mid int64, from, to int64, status int) error {
	time.Sleep(c.s)
	return nil
}

func (c *chatMsgMock) GetChatMessageMidAfter(form, to int64, midAfter int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) GetChatMessageMidSpan(from, to int64, midStart, midEnd int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) AddOfflineMessage(uid int64, mid int64) error {
	time.Sleep(c.s)
	return nil
}

func (c *chatMsgMock) GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	panic("implement me")
}

func (c *chatMsgMock) DelOfflineMessage(uid int64, mid []int64) error {
	panic("implement me")
}

type commMock struct {
	id int64
}

func MockCommDao() {
	Comm = &commMock{}
}

func (c *commMock) GetMessageID() (int64, error) {
	return atomic.AddInt64(&c.id, 1), nil
}
