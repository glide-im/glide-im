package dao

import (
	"errors"
	"go_im/im/entity"
	"go_im/pkg/db"
	"time"
)

func InitMessageDao() {

}

var MessageDao = &messageDao{
	keyMessageIdIncr: "user:message:chat:incr_id",
}

type messageDao struct {
	keyMessageIdIncr string
}

func (m *messageDao) GetUserChatList(uid int64) ([]*Chat, error) {

	var chats []*Chat
	err := db.DB.Table("im_chat").Where("owner = ?", uid).Find(&chats)
	return chats, err.Error
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

	row := 0
	db.DB.Model(&c).Where(Chat{Owner: uid, Target: target, ChatType: typ}).Count(&row)
	if row > 0 {
		return errors.New("chat exist")
	}
	if db.DB.Model(&c).Create(&c).RowsAffected <= 0 {
		return errors.New("create chat error")
	}
	return nil
}

// NewChatMessage
func (m *messageDao) NewChatMessage(sender int64, message *entity.SenderChatMessage) (uint64, error) {

	mid, err := db.Redis.Incr(m.keyMessageIdIncr).Result()
	if err != nil {
		return 0, err
	}

	cm := ChatMessage{
		Mid:         mid,
		Cid:         message.ChatId,
		SenderUid:   sender,
		SendAt:      message.SendAt,
		Message:     message.Message,
		MessageType: message.MessageType,
		At:          "",
	}

	if db.DB.Model(&cm).Create(&cm).RowsAffected <= 0 {
		return 0, errors.New("create chat message error")
	}

	return 0, nil
}

func (m *messageDao) GetChatHistory(cid uint64, uid int64, size int) []*ChatMessage {

	var messages []*ChatMessage

	db.DB.Where("cid = ? and uid = ?", cid, uid).Limit(size).Find(&messages)

	return messages
}

func (m *messageDao) GetChatInfo(chatId uint64) *Chat {

	return &Chat{}
}
