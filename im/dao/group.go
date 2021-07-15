package dao

import "go_im/im"

var GroupDao = new(groupDao)

type groupDao struct{}

func (d *groupDao) GetGroup(gid uint64) *im.Group {

	return im.NewGroup(1, "group", []int64{})
}
