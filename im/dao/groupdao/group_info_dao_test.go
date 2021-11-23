package groupdao

import (
	"go_im/pkg/db"
	"testing"
)

var dao = GroupInfoDaoImpl{}

func init() {
	db.Init()
}

func TestGroupInfoDaoImpl_CreateGroup(t *testing.T) {
	group, err := dao.CreateGroup("MyGroup2", 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(group)
	}
}

func TestGroupInfoDaoImpl_UpdateGroupName(t *testing.T) {
	err := dao.UpdateGroupName(2, "MyGroup")
	t.Log(err)
}

func TestGroupInfoDaoImpl_GetGroupFlag(t *testing.T) {
	flag, err := dao.GetGroupFlag(2)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(flag)
	}
}
