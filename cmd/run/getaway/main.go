package main

import (
	"go_im/service"
	"go_im/service/api_service"
	"go_im/service/gateway"
	"go_im/service/group_messaging"
)

func main() {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = group_messaging.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = api_service.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = gateway.RunServer(config)
	if err != nil {
		panic(err)
	}
}
