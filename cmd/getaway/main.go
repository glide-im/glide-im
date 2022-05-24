package main

import (
	"github.com/glide-im/glideim/im"
	"github.com/glide-im/glideim/service"
	"github.com/glide-im/glideim/service/gateway"
	"github.com/glide-im/glideim/service/messaging_service"
)

func main() {

	im.Init()

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
