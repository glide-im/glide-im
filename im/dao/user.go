package dao

import (
	"errors"
	"go_im/im/model"
	"go_im/pkg/db"
)

var UserDao = new(userDao)

type userDao struct{}

func (d *userDao) GetUser(uid int64) (*model.User, error) {

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
