package msgdao

type groupMsgDao struct {
}

func (groupMsgDao) GetGroupMsgSeq(gid int64) (int64, error) {
	panic("implement me")
}

func (groupMsgDao) UpdateGroupMsgSeq(gid int64, seq int64) error {
	panic("implement me")
}

func (groupMsgDao) GetGroupMessage(mid int64) (*GroupMessageModel, error) {
	panic("implement me")
}

func (groupMsgDao) GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessageModel, error) {
	panic("implement me")
}

func (groupMsgDao) AddGroupMessage(from, to int64, cliMsgID string, type_ int, content string) (*GroupMessageModel, error) {
	panic("implement me")
}

func (groupMsgDao) UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAg int64, lastMsgSeq int64) error {
	panic("implement me")
}

func (groupMsgDao) GetGroupMessageState(gid int64) (*GroupMessageStateModel, error) {
	panic("implement me")
}

func (groupMsgDao) UpdateGroupMemberMsgState(gid int64, uid int64, lastAck int64, lastAckSeq int64) error {
	panic("implement me")
}

func (groupMsgDao) GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgStateModel, error) {
	panic("implement me")
}
