package dao

import (
	"errors"
	"go_im/pkg/db"
	"time"
)

type User struct {
	Uid      int64
	Nickname string
	Avatar   string

	CreateAt time.Time
	UpdateAt time.Time
}

var UserDao = new(userDao)

type userDao struct{}

func (d *userDao) GetUser(uid int64) (*User, error) {

	db.DB.Raw("select * from user where uid = ?", uid)

	return nil, nil
}

func (d *userDao) GetUid(token string) (int64, error) {

	if len(token) == 0 {
		return 0, errors.New("unauthorized")
	}

	// query redis
	return 1, nil
}

func (d *userDao) GetMessageList(uid int64) []uint64 {

	return []uint64{}
}

func (d *userDao) GetFriends(uid int64) []int64 {

	return []int64{}
}
