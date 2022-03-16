package main

import (
	"go_im/im"
	"go_im/service"
	"go_im/service/dispatch"
	"go_im/service/group_messaging"
)

func main() {

	im.Init()

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = dispatch.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = group_messaging.RunServer(config)
	if err != nil {
		panic(err)
	}
}
