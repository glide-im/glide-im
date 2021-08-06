package dao

import (
	"errors"
	"go_im/pkg/db"
	"time"
)

const (
	ChatTypeUser  = 1
	ChatTypeGroup = 2
)

func InitMessageDao() {

}

var ChatDao = &chatDao{
	keyMessageIdIncr:  "user:message:message:incr_id",
	keyChatIdIncr:     "user:message:chat:incr_id",
	keyUserChatIdIncr: "user:message:user_chat:incr_id",
}

type chatDao struct {
	keyMessageIdIncr  string
	keyChatIdIncr     string
	keyUserChatIdIncr string
}

func (m *chatDao) GetChat(target int64, typ int8) (*Chat, error) {

	c := new(Chat)
	err := db.DB.Table("im_chat").Where("target = ? and type = ?", target, typ).Limit(1).Find(c).Error
	return c, err
}

func (m *chatDao) GetUserChatList(uid int64) ([]*UserChat, error) {

	var chats []*UserChat
	err := db.DB.Table("im_user_chat").Where("owner = ?", uid).Find(&chats)
	return chats, err.Error
}

func (m *chatDao) UpdateChatEnterTime(ucId int64) error {
	chat := UserChat{ReadAt: Timestamp(time.Now())}
	db.DB.Model(&chat).Where("uc_id = ?", ucId).Update("read_at")
	return nil
}

func (m *chatDao) UpdateUserChatMsgTime(cid int64, uid int64) (*UserChat, error) {

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

func (m *chatDao) CreateChat(typ int8, targetId int64) (*Chat, error) {

	now := Timestamp(time.Now())
	cid, err := db.Redis.Incr(m.keyChatIdIncr).Result()

	if err != nil {
		return nil, err
	}

	c := Chat{
		Cid:          cid,
		TargetId:     targetId,
		ChatType:     typ,
		CreateAt:     now,
		NewMessageAt: now,
	}

	if db.DB.Model(&c).Create(&c).RowsAffected <= 0 {
		return nil, errors.New("create chat error")
	}
	return &c, nil
}

func (m *chatDao) NewUserChat(cid int64, uid int64, target int64, typ int8) (*UserChat, error) {

	now := nowTimestamp()
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
		ReadAt:   now,
	}

	if db.DB.Model(&uc).Create(&uc).RowsAffected <= 0 {
		return nil, errors.New("create user chat error")
	}

	return &uc, nil
}

// NewChatMessage
func (m *chatDao) NewChatMessage(cid int64, sender int64, msg string, typ int8) (*ChatMessage, error) {

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

func (m *chatDao) GetChatHistory(cid int64, size int) ([]*ChatMessage, error) {

	var messages []*ChatMessage

	err := db.DB.Where("cid = ?", cid).Order("send_at desc").Limit(size).Find(&messages).Error

	return messages, err
}

func (m *chatDao) GetUserChatFromChat(cid int64, uid int64) (*UserChat, error) {

	c := new(UserChat)
	err := db.DB.Table("im_user_chat").Where("cid = ? and owner = ?", cid, uid).Find(c).Error

	return c, err
}
