package dao

import (
	"errors"
	"go_im/pkg/db"
	"time"
)

func InitMessageDao() {

}

var MessageDao = &messageDao{
	keyMessageIdIncr:  "user:message:message:incr_id",
	keyChatIdIncr:     "user:message:chat:incr_id",
	keyUserChatIdIncr: "user:message:user_chat:incr_id",
}

type messageDao struct {
	keyMessageIdIncr  string
	keyChatIdIncr     string
	keyUserChatIdIncr string
}

func (m *messageDao) GetUserChatList(uid int64) ([]*UserChat, error) {

	var chats []*UserChat
	err := db.DB.Table("im_user_chat").Where("owner = ?", uid).Find(&chats)
	return chats, err.Error
}

func (m *messageDao) UpdateChatEnterTime(ucId int64) error {
	chat := UserChat{ReadAt: Timestamp(time.Now())}
	db.DB.Model(&chat).Where("uc_id = ?", ucId).Update("read_at")
	return nil
}

func (m *messageDao) UpdateUserChatMsgTime(cid uint64, uid int64) (*UserChat, error) {

	uc := new(UserChat)
	err := db.DB.Model(uc).Where("cid = ? and owner = ?", cid, uid).Find(&uc).Error
	if err != nil {
		return nil, err
	}
	uc.NewMessageAt = Timestamp(time.Now())
	uc.Unread = uc.Unread + 1
	err = db.DB.Model(uc).Update("new_message_at", "unread").Error

	return uc, err
}

func (m *messageDao) NewChat(uid int64, target uint64, typ int8) (*Chat, error) {

	now := Timestamp(time.Now())
	cid, err := db.Redis.Incr(m.keyChatIdIncr).Result()

	if err != nil {
		return nil, err
	}

	row := 0
	db.DB.Table("im_user_chat").Where("owner = ? and target = ?", uid, target).Count(&row)
	if row > 0 {
		return nil, errors.New("chat exist")
	}

	c := Chat{
		Cid:      cid,
		ChatType: typ,
		CreateAt: now,
	}

	if db.DB.Model(&c).Create(&c).RowsAffected <= 0 {
		return nil, errors.New("create chat error")
	}
	return &c, nil
}

func (m *messageDao) NewUserChat(cid int64, uid int64, target uint64, typ int8) (*UserChat, error) {

	now := Timestamp(time.Now())
	ucid, err := db.Redis.Incr(m.keyUserChatIdIncr).Result()

	if err != nil {
		return nil, err
	}

	uc := UserChat{
		UcId:     ucid,
		Cid:      cid,
		Owner:    uid,
		Target:   target,
		ChatType: typ,
		Unread:   0,
		CreateAt: now,
	}

	if db.DB.Model(&uc).Create(&uc).RowsAffected <= 0 {
		return nil, errors.New("create user chat error")
	}

	return &uc, nil
}

// NewChatMessage
func (m *messageDao) NewChatMessage(cid uint64, sender int64, msg string, typ int8) (*ChatMessage, error) {

	mid, err := db.Redis.Incr(m.keyMessageIdIncr).Result()
	if err != nil {
		return nil, err
	}

	cm := ChatMessage{
		Mid:         mid,
		Cid:         cid,
		SenderUid:   sender,
		SendAt:      Timestamp(time.Now()),
		Message:     msg,
		MessageType: typ,
		At:          "",
	}

	if db.DB.Model(&cm).Create(&cm).RowsAffected <= 0 {
		return nil, errors.New("create chat message error")
	}

	return &cm, nil
}

func (m *messageDao) GetChatHistory(cid uint64, size int) ([]*ChatMessage, error) {

	var messages []*ChatMessage

	err := db.DB.Where("cid = ?", cid).Order("send_at desc").Limit(size).Find(&messages).Error

	return messages, err
}

func (m *messageDao) GetUserChatFromChat(cid uint64, uid int64) (*UserChat, error) {

	c := new(UserChat)
	err := db.DB.Table("im_user_chat").Where("cid = ? and uid = ?", cid, uid).Find(c).Error

	return c, err
}
