package dao

import (
	"go_im/im/dao/uid"
	"go_im/pkg/db"
)

func Init() {

	tables := []interface{}{
		User{},
		Contacts{},
		Chat{},
		ChatMessage{},
		UserChat{},
		Group{},
		GroupMember{},
		GroupMessage{},
	}

	for _, tb := range tables {
		if !db.DB.HasTable(tb) {
			db.DB.CreateTable(tb)
		}
	}

	InitUserDao()
	InitMessageDao()
	uid.Init()
}
