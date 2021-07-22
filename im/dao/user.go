package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"go_im/pkg/db"
	"time"
)

var UserDao *userDao

const userTokenLen = 10
const userTokenExpireDuration = time.Hour * 24 * 3

type userDao struct {
	redisConfig
	mySqlConf
}

func InitUserDao() {
	UserDao = &userDao{
		redisConfig: redisConfig{},
		mySqlConf:   mySqlConf{},
	}

	tb := &User{}
	if !db.DB.HasTable(tb) {
		db.DB.CreateTable(&tb)
	}
}

func (d *userDao) GetUser(uid int64) (*User, error) {

	db.DB.Raw("select * from user where uid = ?", uid)

	return nil, nil
}

func (d *userDao) AddUser(account string, password string) error {

	var count int
	db.DB.Table(d.getUserTableName()).Where("account = ?", account).Select("uid").Count(&count)
	if count > 0 {
		return errors.New("account already exist")
	}

	u := User{
		Account:  account,
		Password: password,
		Avatar:   "",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	if db.DB.Model(&u).Create(&u).RowsAffected > 0 {
		return nil
	} else {
		return errors.New("create account failed")
	}
}

// GetUidByLogin
//
// return uid,token,error
func (d *userDao) GetUidByLogin(account string, password string) (int64, string, error) {

	where := db.DB.Table(d.getUserTableName()).Where("account = ? and password = ?", account, password)
	row := where.Select("uid").Row()

	var uid int64
	if err := row.Scan(&uid); err != nil {
		if err == sql.ErrNoRows {
			return -1, "", errors.New("account does not exist")
		}
		return -1, "", err
	}

	token := genToken(userTokenLen)
	r := db.Redis.Set(d.getKeyUserToken(uid), token, userTokenExpireDuration)
	r2 := db.Redis.Set(token, uid, userTokenExpireDuration)

	if r.Err() != nil || r2.Err() != nil {
		return -1, "", errors.New("redis error")
	}

	return uid, token, nil
}

func (d *userDao) GenToken(uid int64) (string, error) {

	key := d.getKeyUserToken(uid)
	r := db.Redis.Get(key)
	var token string
	if err := r.Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (d *userDao) GetUid(token string) int64 {

	var uid int64
	r := db.Redis.Get(token)
	if err := r.Scan(&uid); err != nil {
		return 0
	}
	return uid
}

func (d *userDao) GetMessageList(uid int64) []uint64 {
	return []uint64{}
}

func (d *userDao) GetFriends(uid int64) []int64 {

	return []int64{}
}

type redisConfig struct{}

func (r *redisConfig) getKeyUserToken(uid int64) string {
	return fmt.Sprintf("user:token:%d", uid)
}

type mySqlConf struct{}

func (m *mySqlConf) getUserTableName() string {
	return "im_user"
}
