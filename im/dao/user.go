package dao

import (
	"errors"
	"fmt"
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

type userDao struct {
	redisConfig
}

func (d *userDao) GetUser(uid int64) (*User, error) {

	db.DB.Raw("select * from user where uid = ?", uid)

	return nil, nil
}

func (d *userDao) GetUidByLogin(username string, password string) (int64, string, error) {

	where := db.DB.Table("im_user").Where("username = ? and password = ?", username, password)
	row := where.Select("uid").Row()

	var uid int64
	if err := row.Scan(&uid); err != nil {
		return -1, "", err
	}

	token := genToken(64)
	r := db.Redis.Set(d.redisConfig.getKeyUserToken(uid), token, time.Hour*24*3)
	if r.Err() != nil {
		return -1, "", errors.New("redisConfig error")
	}

	return uid, token, nil
}

func (d *userDao) GenToken(uid int64) string {

	return ""
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

type redisConfig struct {
}

func (r *redisConfig) getKeyUserToken(uid int64) string {
	return fmt.Sprintf("user:token:%d", uid)
}
