package msgdao

import "go_im/pkg/db"

type chatMsgDao struct {
}

func (chatMsgDao) GetChatMessage(mid int64) (*ChatMessage, error) {
	panic("implement me")
}

func (chatMsgDao) AddOrUpdateChatMessage(message *ChatMessage) (bool, error) {

	c := 0
	err := db.DB.Model(message).Where("m_id = ?", message.MID).Count(&c).Error
	if err != nil {
		return false, err
	}
	if c > 0 {
		return false, nil
	}
	err = db.DB.Model(message).Create(message).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (chatMsgDao) GetChatMessageSeqAfter(uid int64, seqAfter int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (chatMsgDao) GetChatMessageSeqSpan(uid int64, seq int64) (int, error) {
	panic("implement me")
}

func (chatMsgDao) AddOfflineMessage(uid int64, mid int64) error {

	offlineMessage := &OfflineMessage{
		MID: mid,
		UID: uid,
	}
	err := db.DB.Model(offlineMessage).Create(offlineMessage).Error
	return err
}

func (chatMsgDao) GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	panic("implement me")
}

func (chatMsgDao) DelOfflineMessage(uid int64, mid []int64) error {
	panic("implement me")
}
