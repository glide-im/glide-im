package userdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
)

var ContactsDao = &ContactsDaoImpl{}

type ContactsDaoImpl struct{}

func (c ContactsDaoImpl) AddContacts(uid int64, id int64, type_ int8) error {
	contactsID := strconv.FormatInt(uid, 10) + "_" +
		strconv.FormatInt(int64(type_), 10) + "_" +
		strconv.FormatInt(id, 10)
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
	contactsID := strconv.FormatInt(uid, 10) + "_" +
		strconv.FormatInt(int64(type_), 10) + "_" +
		strconv.FormatInt(id, 10)
	query := db.DB.Where("fid = ?", contactsID).Delete(&Contacts{})
	return common.ResolveError(query)
}

func (c ContactsDaoImpl) GetContacts(uid int64) ([]*Contacts, error) {
	var cs []*Contacts
	query := db.DB.Model(&Contacts{}).Where("uid = ?", uid).Find(&cs)
	return cs, common.ResolveError(query)
}
