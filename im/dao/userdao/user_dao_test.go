package userdao

import (
	"go_im/pkg/db"
	"testing"
)

func init() {
	db.Init()
}

func TestUserInfoDaoImpl_AddUser(t *testing.T) {
	err := UserInfoDao.AddUser(&User{
		Account:  "aa",
		Nickname: "nike",
		Password: "",
		Avatar:   "",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUserInfoDaoImpl_UpdatePassword(t *testing.T) {
	err := UserInfoDao.UpdatePassword(543603, "1234567")
	if err != nil {
		t.Error(err)
	}
}

func TestUserInfoDaoImpl_GetUserSimpleInfo(t *testing.T) {
	info, err := UserInfoDao.GetUserSimpleInfo(543602, 543603)
	if err != nil {
		t.Error(err)
	}
	for _, i := range info {
		t.Log(i)
	}
}

func TestUserInfoDaoImpl_GetUserInfo(t *testing.T) {
	info, err := UserInfoDao.GetUser(543603)
	t.Log(info, err)
}
