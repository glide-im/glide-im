package main

import (
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/broker"
)

func main() {

	db.Init()

	configs, err := service.GetConfig()

	if err != nil {
		panic(err)
	}

	err = broker.RunServer(configs)

	if err != nil {
		panic(err)
	}
}
