package main

import (
	"github.com/glide-im/glideim/service"
	"github.com/glide-im/glideim/service/dispatch"
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
