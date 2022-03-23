package main

import (
	"go_im/service"
	"go_im/service/gateway"
	"go_im/service/messaging_service"
)

func main() {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = messaging_service.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = gateway.RunServer(config)
	if err != nil {
		panic(err)
	}
}
