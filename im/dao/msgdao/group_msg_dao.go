package msgdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
	"time"
)

var GroupMsgDaoImpl GroupMsgDao = groupMsgDaoImpl{}

type groupMsgDaoImpl struct {
}

func (groupMsgDaoImpl) GetGroupMsgSeq(gid int64) (int64, error) {
	m := &GroupMsgSeq{}
	query := db.DB.Model(m).Where("gid = ?", gid).Find(m)
	if err := common.ResolveError(query); err != nil {
		return 0, err
	}
	return m.Seq, nil
}

func (groupMsgDaoImpl) UpdateGroupMsgSeq(gid int64, seq int64) error {
	query := db.DB.Model(&GroupMsgSeq{}).Where("gid = ?", gid).Update("seq", seq)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) CreateGroupMsgSeq(gid int64, step int64) error {
	model := &GroupMsgSeq{
		Gid:  gid,
		Seq:  0,
		Step: step,
	}
	query := db.DB.Model(model).Create(model)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) GetLatestGroupMessage(gid int64, pageSize int) ([]*GroupMessage, error) {
	//goland:noinspection GoPreferNilSlice
	ms := []*GroupMessage{}
	query := db.DB.Model(&GroupMessage{}).
		Where("`to` = ?", gid).
		Order("`send_at` DESC").
		Limit(pageSize).
		Find(&ms)
	if err := common.JustError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (groupMsgDaoImpl) GetGroupMessage(gid int64, beforeSeq int64, pageSize int) ([]*GroupMessage, error) {

	//goland:noinspection GoPreferNilSlice
	ms := []*GroupMessage{}
	query := db.DB.Model(&GroupMessage{}).
		Where("`to` = ? AND `seq` < ?", gid, beforeSeq).
		Order("`send_at` DESC").
		Limit(pageSize).
		Find(&ms)
	if err := common.JustError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (groupMsgDaoImpl) GetMessage(mid int64) (*GroupMessage, error) {
	gm := &GroupMessage{}
	query := db.DB.Model(gm).Where("m_id = ?", mid).Find(gm)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return gm, nil
}

func (groupMsgDaoImpl) GetMessages(mid ...int64) ([]*GroupMessage, error) {
	//goland:noinspection GoPreferNilSlice
	gm := []*GroupMessage{}
	query := db.DB.Model(gm).Where("m_id IN (?)", mid).Find(gm)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return gm, nil
}

func (groupMsgDaoImpl) GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error) {
	var ms []*GroupMessage
	query := db.DB.Model(&GroupMessage{}).Where("`to` = ? AND seq > ?", gid, seqAfter).Find(&ms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return ms, nil
}

func (groupMsgDaoImpl) AddGroupMessage(message *GroupMessage) error {
	query := db.DB.Create(message)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) CreateGroupMemberMsgState(gid int64, uid int64) error {
	mbId := strconv.FormatInt(gid, 10) + strconv.FormatInt(uid, 10)
	model := &GroupMemberMsgState{
		MbID:       mbId,
		Gid:        gid,
		UID:        uid,
		LastAckMID: 0,
		LastAckSeq: 0,
	}
	query := db.DB.Create(model)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAt int64, lastMsgSeq int64) error {
	state := &GroupMessageState{
		Gid:       gid,
		LastMID:   lastMID,
		LastSeq:   lastMsgSeq,
		LastMsgAt: lastMsgAt,
	}

	query := db.DB.Model(state).Where("gid = ?", gid).Updates(state)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) CreateGroupMessageState(gid int64) (*GroupMessageState, error) {
	state := &GroupMessageState{
		Gid:       gid,
		LastMsgAt: time.Now().Unix(),
	}
	query := db.DB.Create(state)

	if err := common.ResolveUpdateErr(query); err != nil {
		return nil, err
	}
	return state, nil
}

func (groupMsgDaoImpl) GetGroupMessageState(gid int64) (*GroupMessageState, error) {
	state := &GroupMessageState{}
	query := db.DB.Model(state).Where("gid = ?", gid).Find(state)

	if err := common.ResolveUpdateErr(query); err != nil {
		return nil, err
	}
	return state, nil
}

func (groupMsgDaoImpl) GetGroupsMessageState(gid ...int64) ([]*GroupMessageState, error) {
	//goland:noinspection GoPreferNilSlice
	state := []*GroupMessageState{}
	query := db.DB.Model(state).Where("gid IN (?)", gid).Find(&state)
	if err := common.JustError(query); err != nil {
		return nil, err
	}
	return state, nil
}

func (groupMsgDaoImpl) UpdateGroupMemberMsgState(gid int64, uid int64, ackMid int64, ackSeq int64) error {
	mbId := strconv.FormatInt(gid, 10) + strconv.FormatInt(uid, 10)
	s := &GroupMemberMsgState{}
	query := db.DB.Model(s).Where("mb_id = ?", mbId).Updates(GroupMemberMsgState{
		LastAckMID: ackMid,
		LastAckSeq: ackSeq,
	})
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (groupMsgDaoImpl) GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error) {
	state := &GroupMemberMsgState{}
	mbId := strconv.FormatInt(gid, 10) + strconv.FormatInt(uid, 10)
	query := db.DB.Model(state).Where("mb_id = ?", mbId).Find(state)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return state, nil
}
