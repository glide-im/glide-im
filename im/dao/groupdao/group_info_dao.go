package groupdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"time"
)

type GroupInfoDaoImpl struct {
}

func (GroupInfoDaoImpl) CreateGroup(name string, flag int) (*GroupModel, error) {
	model := &GroupModel{
		Name:     name,
		Mute:     false,
		Flag:     flag,
		CreateAt: time.Now().Unix(),
	}
	query := db.DB.Create(model)
	if err := common.ResolveError(query); err != nil {
		return nil, err
	}
	return model, nil
}

func (GroupInfoDaoImpl) GetGroup(gid int64) (*GroupModel, error) {
	model := &GroupModel{}
	query := db.DB.Model(model).Where("gid = ?", gid).Find(model)
	if err := common.MustFind(query); err != nil {
		return nil, err
	}
	return model, nil
}

func (GroupInfoDaoImpl) GetGroups(gid ...int64) ([]*GroupModel, error) {
	//goland:noinspection GoPreferNilSlice
	model := []*GroupModel{}
	query := db.DB.Model(model).Where("gid IN (?)", gid).Find(&model)
	if err := common.MustFind(query); err != nil {
		return nil, err
	}
	return model, nil
}

func (g *GroupInfoDaoImpl) UpdateGroupName(gid int64, name string) error {
	return g.updateGroupField(gid, "name", name)
}

func (g *GroupInfoDaoImpl) UpdateGroupAvatar(gid int64, avatar string) error {
	return g.updateGroupField(gid, "avatar", avatar)
}

func (g *GroupInfoDaoImpl) UpdateGroupMute(gid int64, mute bool) error {
	return g.updateGroupField(gid, "mute", mute)
}

func (g *GroupInfoDaoImpl) UpdateGroupFlag(gid int64, flag int) error {
	return g.updateGroupField(gid, "flag", flag)
}

func (GroupInfoDaoImpl) GetGroupMute(gid int64) (bool, error) {
	model := &GroupModel{}
	var mute bool
	query := db.DB.Model(model).Where("gid = ?", gid).Select("mute").Limit(1).Find(&mute)
	if err := common.ResolveError(query); err != nil {
		return false, err
	}
	return mute, nil
}

func (GroupInfoDaoImpl) HasGroup(gid int64) (bool, error) {
	var count int64
	query := db.DB.Model(&GroupModel{}).Where("gid = ?", gid).Count(&count)
	if err := common.JustError(query); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (GroupInfoDaoImpl) GetGroupFlag(gid int64) (int, error) {
	model := &GroupModel{}
	var flag int
	query := db.DB.Model(model).Where("gid = ?", gid).Select("flag").Limit(1).Find(&flag)
	if err := common.ResolveError(query); err != nil {
		return flag, err
	}
	return flag, nil
}

func (GroupInfoDaoImpl) updateGroupField(gid int64, field string, value interface{}) error {
	model := &GroupModel{Gid: gid}
	query := db.DB.Model(model).Update(field, value)
	if err := common.ResolveUpdateErr(query); err != nil {
		return err
	}
	return nil
}
