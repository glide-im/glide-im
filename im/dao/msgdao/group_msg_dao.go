package msgdao

import (
	"github.com/pkg/errors"
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

type groupMsgDao struct {
}

func (groupMsgDao) GetGroupMsgSeq(gid int64) (int64, error) {
	m := &GroupMsgSeq{}
	query := db.DB.Model(m).Where("gid = ?", gid).Find(m)
	if err := common.ResolveError(query); err != nil {
		return 0, err
	}
	return m.Seq, nil
}

func (groupMsgDao) UpdateGroupMsgSeq(gid int64) error {
	model := &GroupMsgSeq{
		GID: gid,
	}
	query := db.DB.Model(model).Where("gid = ?", gid).Update(model)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDao) CreateGroupMsgSeq(gid int64, step int64) error {
	model := &GroupMsgSeq{
		GID:  gid,
		Seq:  0,
		Step: step,
	}
	query := db.DB.Model(model).Create(model)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDao) GetGroupMessage(mid int64) (*GroupMessage, error) {
	panic("implement me")
}

func (groupMsgDao) GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error) {
	panic("implement me")
}

func (groupMsgDao) AddGroupMessage(message *GroupMessage) error {
	create := db.DB.Model(message).Create(message)
	if create.Error != nil {
		return create.Error
	}
	if create.RowsAffected == 0 {
		return errors.New("add failed, RowsAffected=0")
	}
	return nil
}

func (groupMsgDao) UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAg int64, lastMsgSeq int64) error {
	panic("implement me")
}

func (groupMsgDao) GetGroupMessageState(gid int64) (*GroupMessageState, error) {
	panic("implement me")
}

func (groupMsgDao) UpdateGroupMemberMsgState(gid int64, uid int64, ackMid int64, ackSeq int64) error {
	mbId := strconv.FormatInt(gid, 10) + strconv.FormatInt(uid, 10)
	s := &GroupMemberMsgState{
		MbID:       mbId,
		LastAckMID: ackMid,
		LastAckSeq: ackSeq,
	}
	query := db.DB.Model(s).Where("mb_id = ?", mbId).Update(s)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDao) GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error) {
	panic("implement me")
}
