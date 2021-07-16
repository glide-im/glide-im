package dao

import "go_im/im"

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

func (d *groupDao) GetGroup(gid uint64) *im.Group {

	return im.NewGroup(1, "group", []int64{})
}
