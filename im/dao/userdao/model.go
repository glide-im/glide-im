package userdao

import (
	"go_im/im/dao/common"
)

type User struct {
	Uid      int64  `gorm:"primary_key"`
	Account  string `gorm:"unique"`
	Nickname string
	Password string
	Avatar   string

	CreateAt common.Timestamp `gorm:"type:datetime"`
	UpdateAt common.Timestamp `gorm:"type:datetime"`
}

type Contacts struct {
	Fid      int64 `gorm:"primary_key"`
	Owner    int64
	TargetId int64
	Remark   string
	Type     int8
	AddTime  common.Timestamp `gorm:"type:datetime"`
}
