package dao

import (
	"go_im/im/dao/groupdao"
	"go_im/im/dao/msgdao"
	"go_im/im/dao/uid"
	"go_im/im/dao/userdao"
	"go_im/pkg/db"
)

func Init() {

	tables := []interface{}{
		msgdao.ChatMessage{},
		groupdao.Group{},
		groupdao.GroupMember{},
		groupdao.GroupMessage{},
	}

	for _, tb := range tables {
		if !db.DB.HasTable(tb) {
			db.DB.CreateTable(tb)
		}
	}

	userdao.InitUserDao()
	uid.Init()
}
