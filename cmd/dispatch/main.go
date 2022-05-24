package main

import (
	"go_im/service"
	"go_im/service/dispatch"
)

func main() {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = dispatch.RunServer(config)
	if err != nil {
		panic(err)
	}
}
