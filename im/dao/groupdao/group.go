package groupdao

import (
	"errors"
	"fmt"
	"go_im/im/dao/common"
	"go_im/pkg/db"
)

var GroupDao2 = new(groupDao)

type groupDao struct{}

func (d *groupDao) CreateGroup(name string, owner int64) (*Group, error) {

	gid, err := GetNextGid()
	if err != nil {
		return nil, err
	}

	g := Group{
		Gid:      gid,
		Name:     name,
		Avatar:   "",
		Owner:    owner,
		Mute:     false,
		Notice:   "",
		ChatId:   0,
		CreateAt: common.NowTimestamp(),
	}

	if db.DB.Model(&g).Create(&g).RowsAffected <= 0 {
		return nil, errors.New("create group error")
	}

	return &g, nil
}

func (d *groupDao) GetMember(gid int64, uid ...int64) ([]*GroupMember, error) {

	q := db.DB.Table("im_group_member")
	q = q.Where("gid = ? AND uid IN (?)", gid, uid)

	var mbs []*GroupMember
	err := q.Select("uid").Find(&mbs).Error
	return mbs, err
}

func (d *groupDao) GetAllGroup() ([]*Group, error) {
	var groups []*Group
	err := db.DB.Table("im_group").Find(&groups).Error
	return groups, err
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

func (d *groupDao) AddMember(gid int64, typ int32, uid ...int64) ([]*GroupMember, error) {

	var members []*GroupMember

	for _, i := range uid {
		gm := GroupMember{
			Gid:    gid,
			Uid:    i,
			Mute:   0,
			Remark: "",
			Flag:   typ,
			JoinAt: common.NowTimestamp(),
		}
		if db.DB.Model(&gm).Create(&gm).RowsAffected <= 0 {
			return members, errors.New("add member error")
		}
		members = append(members, &gm)
	}

	if len(members) != len(uid) {
		return nil, errors.New(fmt.Sprintf("add member error, expect size: %d, actully: %d", len(uid), len(members)))
	}
	return members, nil
}

func (d *groupDao) GetMembers(gid int64) ([]*GroupMember, error) {

	var gm []*GroupMember

	err := db.DB.Table("im_group_member").Where("gid = ?", gid).Find(&gm).Error

	return gm, err
}

func (d *groupDao) GetUserGroup(uid int64) ([]*Group, error) {
	var groups []*Group
	err := db.DB.Table("im_group_member").Where("uid = ?", uid).Find(&groups).Error
	return groups, err
}
