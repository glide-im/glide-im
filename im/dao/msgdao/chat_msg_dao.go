package msgdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

type chatMsgDao struct {
}

func (chatMsgDao) GetChatMessage(mid int64) (*ChatMessage, error) {
	m := &ChatMessage{}
	query := db.DB.Model(m).Where("m_id = ?", mid).Find(m)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return m, nil
}

func (chatMsgDao) AddOrUpdateChatMessage(message *ChatMessage) (bool, error) {
	var c int64
	query := db.DB.Table("im_chat_message").Where("m_id = ?", message.MID).Count(&c)
	if err := common.ResolveError(query); err != nil {
		return false, err
	}
	if c > 0 {
		return false, nil
	}
	query = db.DB.Create(message)
	if err := common.ResolveError(query); err != nil {
		return false, err
	}
	return true, nil
}

func (chatMsgDao) GetChatMessageMidAfter(from, to int64, midAfter int64) ([]*ChatMessage, error) {
	lg, sm := from, to
	if lg < sm {
		lg, sm = sm, lg
	}
	sessionTag := strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).Where("session_tag = ? and m_id > ?", sessionTag, midAfter).Find(&ms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (chatMsgDao) GetChatMessageMidSpan(from, to int64, midStart, midEnd int64) ([]*ChatMessage, error) {
	lg, sm := from, to
	if lg < sm {
		lg, sm = sm, lg
	}
	sessionTag := strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).Where("session_tag = ? AND m_id >= ? AND m_id < ?", sessionTag, midStart, midEnd).Find(&ms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (chatMsgDao) AddOfflineMessage(uid int64, mid int64) error {
	offlineMessage := &OfflineMessage{
		MID: mid,
		UID: uid,
	}
	query := db.DB.Create(offlineMessage)
	return common.ResolveError(query)
}

func (chatMsgDao) GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	var m []*OfflineMessage
	query := db.DB.Model(&OfflineMessage{}).Where("uid = ?", uid).Find(&m)
	if query.Error != nil {
		return nil, query.Error
	}
	return m, nil
}

func (chatMsgDao) DelOfflineMessage(uid int64, mid []int64) error {
	query := db.DB.Where("uid = ? AND m_id IN (?)", uid, mid).Delete(&OfflineMessage{})
	return query.Error
}
