package dao

import "time"

type Chat struct {
	Id       uint64
	Owner    int64
	Target   int64
	Unread   int
	ReadAt   time.Time
	CreateAt time.Time
}

type ChatMessage struct {
	Sender      int64
	SendAt      time.Time
	Message     string
	MessageType int
}

var MessageDao = new(messageDao)

type messageDao struct{}

// NewChatMessage
func (m *messageDao) NewChatMessage(chatId uint64, content string, msgType int) error {

	return nil
}

func (m *messageDao) GetChatMessage(chatId uint64, size int) []*ChatMessage {

	return []*ChatMessage{}
}

func (m *messageDao) GetChatInfo(chatId uint64) *Chat {

	return &Chat{}
}
