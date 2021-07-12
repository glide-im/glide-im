package model

import "time"

type User struct {
	Uid      int64
	Nickname string
	Avatar   string

	CreateAt time.Time
	UpdateAt time.Time
}
