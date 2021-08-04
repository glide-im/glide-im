package dao

import "go_im/pkg/db"

func Init() {

	tables := []interface{}{
		User{},
		Friend{},
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
}
