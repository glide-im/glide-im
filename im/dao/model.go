package dao

import "time"

type User struct {
	Uid      int64 `gorm:"primary_key"`
	Account  string
	Password string
	Avatar   string

	CreateAt time.Time
	UpdateAt time.Time
}
