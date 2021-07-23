package dao

import (
	"errors"
	"go_im/pkg/db"
	"time"
)

var MessageDao = new(messageDao)

type messageDao struct{}

func (m *messageDao) GetUserChatList(uid int64) ([]*Chat, error) {

	var chats []*Chat
	err := db.DB.Where("uid=?", uid).Find(chats).Error
	return chats, err
}

func (m *messageDao) NewChat(uid int64, target uint64, typ int8) error {

	now := time.Now()

	c := Chat{
		Owner:        uid,
		Target:       target,
		ChatType:     typ,
		NewMessageAt: now,
		ReadAt:       now,
		CreateAt:     now,
	}

	if db.DB.Model(&c).Create(&c).RowsAffected <= 0 {
		return errors.New("create chat error")
	}
	return nil
}

// NewChatMessage
func (m *messageDao) NewChatMessage(cid uint64, sender int64, content string, at string, msgType int) error {

	cm := ChatMessage{
		Cid:         cid,
		SenderUid:   sender,
		SendAt:      time.Now(),
		Message:     content,
		MessageType: msgType,
		At:          at,
	}

	if db.DB.Model(&cm).Create(&cm).RowsAffected <= 0 {
		return errors.New("create chat message error")
	}

	return nil
}

func (m *messageDao) GetChatHistory(cid uint64, uid int64, size int) []*ChatMessage {

	var messages []*ChatMessage

	db.DB.Where("cid = ? and uid = ?", cid, uid).Limit(size).Find(&messages)

	return messages
}

func (m *messageDao) GetChatInfo(chatId uint64) *Chat {

	return &Chat{}
}
