package main

import (
	"go_im/im/dao"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/api_service"
	"go_im/service/broker"
	"go_im/service/dispatch"
)

func main() {

	db.Init()
	dao.Init()

	config, err := service.GetConfig()
	must(err)

	err = dispatch.SetupClient(config)
	must(err)

	err = broker.SetupClient(config)
	must(err)

	err = api_service.RunServer(config)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
