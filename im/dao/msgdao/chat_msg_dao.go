package msgdao

type chatMsgDao struct {
}

func (chatMsgDao) GetChatMessage(mid int64) (*ChatMessage, error) {
	panic("implement me")
}

func (chatMsgDao) AddChatMessage(from, to int64, cliMsgID string, type_ int, content string) (int64, error) {
	panic("implement me")
}

func (chatMsgDao) GetChatMessageSeqAfter(uid int64, seqAfter int64) ([]*ChatMessage, error) {
	panic("implement me")
}

func (chatMsgDao) GetChatMessageSeqSpan(uid int64, seq int64) (int, error) {
	panic("implement me")
}

func (chatMsgDao) AddOfflineMessage(uid int64, mid int64, seq int64) error {
	panic("implement me")
}

func (chatMsgDao) GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	panic("implement me")
}

func (chatMsgDao) DelOfflineMessage(uid int64, mid []int64) error {
	panic("implement me")
}
