package dao

import (
	"errors"
	"go_im/pkg/db"
)

/**
Member

Type 1: 群员 2: 管理 3: 群主
State 状态位 0000 : 0-0-通知开关-被禁言
*/
type Member struct {
	Uid      int64
	Nickname string
	Avatar   string
	Type     uint8
	State    uint8
}

var GroupDao = new(groupDao)

type groupDao struct{}

func (d *groupDao) NewGroup(name string, owner int64) (*Group, error) {

	g := Group{
		Gid:      0,
		Name:     name,
		Avatar:   "",
		Owner:    owner,
		Mute:     false,
		Notice:   "",
		CreateAt: nowTimestamp(),
	}

	if db.DB.Model(&g).Create(&g).RowsAffected <= 0 {
		return nil, errors.New("create group error")
	}

	return &g, nil
}

func (d *groupDao) GetGroup(gid int64) (*Group, error) {

	g := new(Group)
	err := db.DB.Model(g).Where("gid = ?", gid).Find(g).Error
	return g, err
}

func (d *groupDao) RemoveMember(gid int64, uid int64) error {

	e := db.DB.Table("im_group_member").Delete("gid = ? and uid = ?", gid, uid).Error

	if e != nil {
		return e
	}

	return nil
}

func (d *groupDao) AddMember(gid int64, uid int64, typ int8) error {

	gm := GroupMember{
		Gid:    gid,
		Uid:    uid,
		Mute:   0,
		Remark: "",
		Type:   typ,
		JoinAt: nowTimestamp(),
	}

	if db.DB.Model(&gm).Create(&gm).RowsAffected <= 0 {
		return errors.New("add member error")
	}

	return nil
}

func (d *groupDao) GetMembers(gid int64) ([]*GroupMember, error) {

	var gm []*GroupMember

	err := db.DB.Model(gm).Where("gid = ?", gid).Find(gm).Error

	return gm, err
}

func (d *groupDao) GetUserGroup(uid int64) ([]*Group, error) {
	var groups []*Group
	err := db.DB.Table("im_group_member").Where("uid = ?", uid).Find(&groups).Error
	return groups, err
}
