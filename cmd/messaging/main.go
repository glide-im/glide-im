package main

import (
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/broker"
	"go_im/service/dispatch"
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

	err = broker.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = messaging_service.RunServer(config)

	if err != nil {
		panic(err)
	}
}
