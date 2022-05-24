package im

import (
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/dao"
	"github.com/glide-im/glideim/im/messaging"
	"github.com/glide-im/glideim/pkg/db"
)

func Init() {
	db.Init()
	dao.Init()

	client.SetMessageHandler(messaging.HandleMessage)
}
