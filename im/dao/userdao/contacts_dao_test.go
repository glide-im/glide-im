package userdao

import (
	"go_im/pkg/db"
	"testing"
)

func init() {
	db.Init()
}

func TestContactsDaoImpl_AddContacts(t *testing.T) {
	err := ContactsDao.AddContacts(1, 4, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestContactsDaoImpl_DelContacts(t *testing.T) {
	err := ContactsDao.DelContacts(1, 3, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestContactsDaoImpl_GetContacts(t *testing.T) {
	contacts, err := ContactsDao.GetContacts(1)
	if err != nil {
		t.Error(err)
	}
	for _, c := range contacts {
		t.Log(c)
	}
}
