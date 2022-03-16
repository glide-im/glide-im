package main

import (
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/dispatch"
	"go_im/service/gateway"
	"go_im/service/group_messaging"
	"go_im/service/messaging_service"
)

func main() {
	db.Init()
	dao.Init()

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = dispatch.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = gateway.InitMessageProducer(config.Nsq.Nsqd)
	if err != nil {
		panic(err)
	}

	err = group_messaging.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = messaging_service.RunServer(config)

	if err != nil {
		panic(err)
	}
}
