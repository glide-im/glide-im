package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"go_im/im/dao/uid"
	"go_im/pkg/db"
	"math/rand"
	"time"
)

var UserDao *userDao

const userTokenLen = 10
const userTokenExpireDuration = time.Hour * 24 * 3

const (
	ContactsTypeUser  = 1
	ContactsTypeGroup = 2
)

var avatars = []string{
	"https://dengzii.com/static/a.webp",
	"https://dengzii.com/static/b.webp",
	"https://dengzii.com/static/c.webp",
	"https://dengzii.com/static/d.webp",
	"https://dengzii.com/static/e.webp",
	"https://dengzii.com/static/f.webp",
	"https://dengzii.com/static/g.webp",
	"https://dengzii.com/static/h.webp",
	"https://dengzii.com/static/i.webp",
	"https://dengzii.com/static/j.webp",
	"https://dengzii.com/static/k.webp",
	"https://dengzii.com/static/l.webp",
	"https://dengzii.com/static/m.webp",
	"https://dengzii.com/static/n.webp",
	"https://dengzii.com/static/o.webp",
	"https://dengzii.com/static/p.webp",
	"https://dengzii.com/static/q.webp",
	"https://dengzii.com/static/r.webp",
}
var nextAvatar = 0

var nickName = []string{"佐菲", "赛文", "杰克", "艾斯", "泰罗", "雷欧", "阿斯特拉", "艾迪", "迪迦", "杰斯", "奈克斯", "梦比优斯", "盖亚", "戴拿"}
var nextNickName = 0

type userDao struct {
	redisConfig
	mySqlConf
}

func InitUserDao() {
	UserDao = &userDao{
		redisConfig: redisConfig{},
		mySqlConf:   mySqlConf{},
	}
}

func (d *userDao) HasUser(uid ...int64) (bool, error) {

	query := db.DB.Table("im_user").Where("uid = ?", uid[0])

	for index, id := range uid {
		if index == 0 {
			continue
		}
		query = query.Or("uid = ?", id)
	}
	rows := -1
	err := query.Count(&rows).Error

	return rows == len(uid), err
}

func (d *userDao) GetUser(uid ...int64) ([]*User, error) {

	var u []*User
	query := db.DB.Where("uid = ?", uid[0])

	for index, id := range uid {
		if index == 0 {
			continue
		}
		query = query.Or("uid = ?", id)
	}

	return u, query.Find(&u).Error
}

func (d *userDao) AddUser(account string, password string) error {

	var count int
	db.DB.Table(d.getUserTableName()).Where("account = ?", account).Select("uid").Count(&count)
	if count > 0 {
		return errors.New("account already exist")
	}

	t := Timestamp(time.Now())
	u := User{
		Uid:      uid.GenUid(),
		Account:  account,
		Password: password,
		Nickname: nickName[rand.Intn(14)],
		Avatar:   avatars[rand.Intn(17)],
		CreateAt: t,
		UpdateAt: t,
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

func (d *userDao) GetAllContacts(uid int64) ([]*Contacts, error) {

	var ret []*Contacts
	err := db.DB.Table("im_contacts").Where("owner = ?", uid).Find(&ret).Error
	return ret, err
}

func (d *userDao) HasContacts(owner int64, targetId int64, typ int8) (bool, error) {
	row := 0
	err := db.DB.Table("im_contacts").Where("target_id = ? and owner = ? and type = ?", targetId, owner, typ).Count(&row).Error
	return row > 0, err
}

func (d *userDao) AddContacts(uid int64, targetId int64, typ int8, remark string) (*Contacts, error) {

	f := &Contacts{
		Owner:    uid,
		TargetId: targetId,
		Remark:   remark,
		Type:     typ,
		AddTime:  nowTimestamp(),
	}
	if db.DB.Model(f).Create(f).RowsAffected <= 0 {
		return nil, errors.New("create friend error")
	}
	return f, nil
}

type redisConfig struct{}

func (r *redisConfig) getKeyUserToken(uid int64) string {
	return fmt.Sprintf("user:token:%d", uid)
}

type mySqlConf struct{}

func (m *mySqlConf) getUserTableName() string {
	return "im_user"
}
