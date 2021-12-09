package groupdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

type GroupMemberDaoImpl struct {
}

func (GroupMemberDaoImpl) GetMembers(gid int64) ([]*GroupMemberModel, error) {
	var gms []*GroupMemberModel
	query := db.DB.Model(&GroupMemberModel{}).Where("gid = ?", gid).Find(&gms)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return gms, nil
}

func GetMemberID(gid int64, uid int64) string {
	return strconv.FormatInt(gid, 10) + "_" + strconv.FormatInt(uid, 10)
}

func (GroupMemberDaoImpl) HasMember(gid int64, uid int64) (bool, error) {
	mbID := GetMemberID(gid, uid)
	var count int64
	query := db.DB.Model(&GroupMemberModel{}).Where("mb_id = ?", mbID).Count(&count)
	if err := common.JustError(query); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (GroupMemberDaoImpl) AddMembers(gid int64, flag int64, typ int64, uid ...int64) error {
	var ms []*GroupMemberModel
	for _, i := range uid {
		ms = append(ms, &GroupMemberModel{
			MbID: GetMemberID(gid, i),
			Gid:  gid,
			Uid:  i,
			Flag: flag,
			Type: typ,
		})
	}
	query := db.DB.Create(&ms)
	if err := common.MustUpdate(query); err != nil {
		return err
	}
	return nil
}

func (GroupMemberDaoImpl) AddMember(gid int64, uid int64, typ int64, defaultFlag int64) error {
	mbId := GetMemberID(gid, uid)
	model := GroupMemberModel{
		MbID: mbId,
		Gid:  gid,
		Uid:  uid,
		Flag: defaultFlag,
		Type: typ,
	}
	query := db.DB.Create(model)
	if err := common.ResolveError(query); err != nil {
		return err
	}
	return nil
}

func (GroupMemberDaoImpl) RemoveMember(gid int64, uid int64) error {
	query := db.DB.Where("gid = ? AND uid = ?", gid, uid).Delete(&GroupMember{})
	return common.ResolveError(query)
}

func (GroupMemberDaoImpl) GetMemberFlag(gid int64, uid int64) (int64, error) {
	mbId := GetMemberID(gid, uid)
	var flag int64
	query := db.DB.Model(&GroupMemberModel{}).Where("mb_id = ?", mbId).Select("flag").Find(&flag)
	if err := common.ResolveError(query); err != nil {
		return 0, err
	}
	return flag, nil
}

func (GroupMemberDaoImpl) UpdateMemberFlag(gid int64, uid int64, flag int) error {
	mbId := GetMemberID(gid, uid)
	query := db.DB.Model(&GroupMemberModel{}).Where("mb_id = ?", mbId).Update("flag", flag)
	return common.ResolveError(query)
}

func (GroupMemberDaoImpl) GetMember(gid int64, uid int64) (*GroupMemberModel, error) {
	var gm *GroupMemberModel
	mbId := GetMemberID(gid, uid)
	query := db.DB.Model(&GroupMemberModel{}).Where("mb_id = ?", mbId).Find(gm)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return gm, nil
}
