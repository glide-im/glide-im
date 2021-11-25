package userdao

import (
	"errors"
	"go_im/pkg/db"
)

var UserDao2 *userDao

const userTokenLen = 10

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

var nickName = []string{"佐菲", "赛文", "杰克", "艾斯", "泰罗", "雷欧", "阿斯特拉", "艾迪", "迪迦", "杰斯", "奈克斯", "梦比优斯", "盖亚", "戴拿"}

type userDao struct {
	mySqlConf
}

func InitUserDao() {
	UserDao2 = &userDao{
		mySqlConf: mySqlConf{},
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
	var rows int64
	err := query.Count(&rows).Error

	return rows == int64(len(uid)), err
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

func (d *userDao) GetAllContacts(uid int64) ([]*Contacts, error) {

	var ret []*Contacts
	err := db.DB.Table("im_contacts").Where("owner = ?", uid).Find(&ret).Error
	return ret, err
}

func (d *userDao) HasContacts(owner int64, targetId int64, typ int8) (bool, error) {
	var row int64
	err := db.DB.Table("im_contacts").Where("target_id = ? and owner = ? and type = ?", targetId, owner, typ).Count(&row).Error
	return row > 0, err
}

func (d *userDao) AddContacts(uid int64, targetId int64, typ int8, remark string) (*Contacts, error) {

	f := &Contacts{
		Uid:    uid,
		Id:     targetId,
		Remark: remark,
		Type:   typ,
	}
	if db.DB.Model(f).Create(f).RowsAffected <= 0 {
		return nil, errors.New("create friend error")
	}
	return f, nil
}

type mySqlConf struct{}

func (m *mySqlConf) getUserTableName() string {
	return "im_user"
}
