package msgdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

var ChatMsgDaoImpl ChatMsgDao = chatMsgDaoImpl{}

type chatMsgDaoImpl struct {
}

func (chatMsgDaoImpl) GetRecentChatMessagesBySession(uid1, uid2 int64, pageSize int) ([]*ChatMessage, error) {
	sid, _, _ := getSessionId(uid2, uid1)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).
		Where("`session_id` = ?", sid).
		Order("`send_at` DESC").
		Limit(pageSize).
		Find(&ms)
	if query.Error != nil {
		return nil, query.Error
	}
	return ms, nil
}

func (chatMsgDaoImpl) GetChatMessagesBySession(uid1, uid2 int64, beforeMid int64, pageSize int) ([]*ChatMessage, error) {
	sid, _, _ := getSessionId(uid2, uid1)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).
		Where("`session_id` = ? AND `m_id` < ?", sid, beforeMid).
		Order("`send_at` DESC").
		Limit(pageSize).
		Find(&ms)
	if query.Error != nil {
		return nil, query.Error
	}
	return ms, nil
}

func (chatMsgDaoImpl) GetRecentChatMessages(uid int64, after int64) ([]*ChatMessage, error) {
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).Where("`from` = ? OR `to` = ? AND `send_at` > ?", uid, uid, after).Find(&ms)
	if query.Error != nil {
		return nil, query.Error
	}
	return ms, nil
}

func (chatMsgDaoImpl) GetChatMessage(mid ...int64) ([]*ChatMessage, error) {
	//goland:noinspection GoPreferNilSlice
	m := []*ChatMessage{}
	query := db.DB.Model(m).Where("m_id in (?)", mid).Find(&m)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return m, nil
}

func (chatMsgDaoImpl) AddOrUpdateChatMessage(message *ChatMessage) (bool, error) {
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

func (chatMsgDaoImpl) GetChatMessageMidAfter(from, to int64, midAfter int64) ([]*ChatMessage, error) {
	lg, sm := from, to
	if lg < sm {
		lg, sm = sm, lg
	}
	sid := strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).Where("session_id = ? and m_id > ?", sid, midAfter).Find(&ms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (chatMsgDaoImpl) GetChatMessageMidSpan(from, to int64, midStart, midEnd int64) ([]*ChatMessage, error) {
	lg, sm := from, to
	if lg < sm {
		lg, sm = sm, lg
	}
	sid := strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
	var ms []*ChatMessage
	query := db.DB.Model(&ChatMessage{}).Where("sid = ? AND m_id >= ? AND m_id < ?", sid, midStart, midEnd).Find(&ms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (chatMsgDaoImpl) AddOfflineMessage(uid int64, mid int64) error {
	offlineMessage := &OfflineMessage{
		MID: mid,
		UID: uid,
	}
	query := db.DB.Create(offlineMessage)
	return common.ResolveError(query)
}

func (chatMsgDaoImpl) GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	var m []*OfflineMessage
	query := db.DB.Model(&OfflineMessage{}).Where("uid = ?", uid).Find(&m)
	if query.Error != nil {
		return nil, query.Error
	}
	return m, nil
}

func (chatMsgDaoImpl) DelOfflineMessage(uid int64, mid []int64) error {
	query := db.DB.Where("uid = ? AND m_id IN (?)", uid, mid).Delete(&OfflineMessage{})
	return query.Error
}
