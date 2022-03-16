package im

import (
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/messaging"
	"go_im/pkg/db"
)

func Init() {
	db.Init()
	dao.Init()

	client.SetMessageHandler(messaging.HandleMessage)
}
