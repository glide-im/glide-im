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

var ChatDao = &chatDao{}

type chatDao struct{}

func (m *chatDao) GetChatByTarget(target int64, typ int8) (*Chat, error) {

	c := new(Chat)
	err := db.DB.Table("im_chat").Where("target_id = ? and chat_type = ?", target, typ).Limit(1).Find(c).Error
	return c, err
}

func (m *chatDao) GetChat(chatId int64) (*Chat, error) {
	c := new(Chat)
	err := db.DB.Table("im_chat").Where("cid = ?", chatId).Limit(1).Find(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (m *chatDao) GetCurrentMessageID(chatId int64) (int64, error) {
	cmid := new(ChatMessageID)
	res := db.DB.Table("im_chat_message_id").Where("cid = ?", chatId).Limit(1).Find(cmid)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("no such chat")
	}
	return cmid.CurrentMid, nil
}

func (m *chatDao) UpdateCurrentMessageID(chatID int64, mid int64) error {
	res := db.DB.
		Table("chat_message_id").
		Where("cid = ?", chatID).
		Update(map[string]interface{}{"current_mid": mid})
	return resolveError(res)
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
	cid, err := GetNextChatId(typ)

	if err != nil {
		return nil, err
	}

	c := Chat{
		Cid:          cid,
		CurrentMid:   1,
		TargetId:     targetId,
		ChatType:     typ,
		CreateAt:     now,
		NewMessageAt: now,
	}

	chatMid := ChatMessageID{
		Cid:        cid,
		CurrentMid: 1,
	}
	db.DB.Model(&chatMid).Create(&chatMid)

	if db.DB.Model(&c).Create(&c).RowsAffected <= 0 {
		return nil, errors.New("create chat error")
	}
	return &c, nil
}

func (m *chatDao) NewUserChat(cid int64, uid int64, target int64, typ int8) (*UserChat, error) {

	now := nowTimestamp()
	ucid, err := GetUserChatId(uid, cid)

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

// NewChatMessage 插入一条新消息到数据库
func (m *chatDao) NewChatMessage(cid int64, sender int64, msg string, typ int8) (*ChatMessage, error) {

	mid := GetNextMessageId(cid)

	cm := ChatMessage{
		Mid:         mid,
		Cid:         cid,
		Sender:      sender,
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

func (m *chatDao) NewGroupMessage(message ChatMessage) error {
	err := db.DB.Model(&message).Create(&message).Error
	return err
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
