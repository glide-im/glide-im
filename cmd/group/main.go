package main

import (
	"github.com/glide-im/glideim/im"
	"github.com/glide-im/glideim/service"
	"github.com/glide-im/glideim/service/dispatch"
	"github.com/glide-im/glideim/service/group_messaging"
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
