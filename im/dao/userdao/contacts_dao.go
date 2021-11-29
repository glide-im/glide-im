package userdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

var ContactsDao = &ContactsDaoImpl{}

type ContactsDaoImpl struct{}

func getContactsId(uid int64, id int64, type_ int8) string {
	return strconv.FormatInt(uid, 10) + "_" +
		strconv.FormatInt(int64(type_), 10) + "_" +
		strconv.FormatInt(id, 10)
}

func (c ContactsDaoImpl) HasContacts(uid int64, id int64, type_ int8) (bool, error) {
	contactsID := getContactsId(uid, id, type_)
	var count int64
	query := db.DB.Model(&Contacts{}).Where("f_id = ?", contactsID).Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

func (c ContactsDaoImpl) AddContacts(uid int64, id int64, type_ int8) error {
	contactsID := getContactsId(uid, id, type_)
	contacts := &Contacts{
		Fid:    contactsID,
		Uid:    uid,
		Id:     id,
		Remark: "",
		Type:   type_,
	}
	query := db.DB.Create(contacts)
	return common.ResolveError(query)
}

func (c ContactsDaoImpl) DelContacts(uid int64, id int64, type_ int8) error {
	contactsID := getContactsId(uid, id, type_)
	query := db.DB.Where("fid = ?", contactsID).Delete(&Contacts{})
	return common.ResolveError(query)
}

func (c ContactsDaoImpl) GetContacts(uid int64) ([]*Contacts, error) {
	var cs []*Contacts
	query := db.DB.Model(&Contacts{}).Where("uid = ?", uid).Find(&cs)
	return cs, common.ResolveError(query)
}
