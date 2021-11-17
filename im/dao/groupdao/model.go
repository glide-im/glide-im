package groupdao

import "go_im/im/dao/common"

type Group struct {
	Gid      int64 `gorm:"primary_key"`
	Name     string
	Avatar   string
	Owner    int64
	Mute     bool
	Notice   string
	ChatId   int64
	CreateAt common.Timestamp `gorm:"type:datetime"`
}

type GroupMember struct {
	Id     int64 `gorm:"primary_key"`
	Gid    int64
	Uid    int64
	Mute   int64
	Remark string
	Flag   int32
	JoinAt common.Timestamp `gorm:"type:datetime"`
}

type GroupMessage struct {
	GmId        int64 `gorm:"primary_key"`
	Cid         int64
	SenderUid   int64
	SendAt      common.Timestamp `gorm:"type:datetime"`
	Message     string
	MessageType int8
	At          string
}

type GroupNotify struct {
	Id     int64 `gorm:"primary_key"`
	Gid    int64
	Uid    int64
	Remark string
	Type   int8
	State  int8
}
