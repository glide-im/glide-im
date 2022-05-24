package main

import (
	"github.com/glide-im/glideim/pkg/db"
	"github.com/glide-im/glideim/service"
	"github.com/glide-im/glideim/service/broker"
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
