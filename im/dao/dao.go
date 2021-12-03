package dao

import (
	"go_im/im/dao/uid"
	"go_im/im/dao/userdao"
)

func Init() {
	userdao.InitUserDao()
	uid.Init()
}
