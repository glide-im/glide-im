package groupdao

import (
	"go_im/pkg/db"
	"testing"
)

var gmDao = GroupMemberDaoImpl{}

func init() {
	db.Init()
}

func TestGroupMemberDaoImpl_AddMember(t *testing.T) {
	err := gmDao.AddMember(14, 4, 2, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestGroupMemberDaoImpl_GetMemberType(t *testing.T) {
	tp, err := gmDao.GetMemberType(17, 543602)
	if err != nil {
		t.Error(err)
	}
	t.Log(tp)
}

func TestGroupMemberDaoImpl_GetMembers(t *testing.T) {
	mbs, err := gmDao.GetMembers(16)
	if err != nil {
		t.Error(err)
	}
	for _, mb := range mbs {
		t.Log(mb)
	}
}

func TestGroupMemberDaoImpl_UpdateMemberFlag(t *testing.T) {
	err := gmDao.UpdateMemberFlag(1, 3, 4)
	if err != nil {
		t.Error(err)
	}
}

func TestGroupMemberDaoImpl_GetMemberFlag(t *testing.T) {
	flag, err := gmDao.GetMemberFlag(1, 3)
	if err != nil {
		t.Error(err)
	}
	t.Log(flag)
}
