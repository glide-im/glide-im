package main

import (
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/api_service"
	"go_im/service/gateway"
	"go_im/service/group_messaging"
)

func main() {

	config, err := service.GetConfig()
	must(err)

	db.Init()
	dao.Init()

	err = gateway.InitMessageProducer(config.Nsq.Nsqd)
	must(err)

	err = group_messaging.SetupClient(config)
	must(err)

	err = api_service.RunServer(config)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
