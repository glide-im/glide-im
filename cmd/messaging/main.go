package main

import (
	"github.com/glide-im/glideim/im/dao"
	"github.com/glide-im/glideim/pkg/db"
	"github.com/glide-im/glideim/service"
	"github.com/glide-im/glideim/service/broker"
	"github.com/glide-im/glideim/service/dispatch"
	"github.com/glide-im/glideim/service/messaging_service"
)

func main() {
	db.Init()
	dao.Init()

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}

	err = dispatch.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = broker.SetupClient(config)
	if err != nil {
		panic(err)
	}

	err = messaging_service.RunServer(config)

	if err != nil {
		panic(err)
	}
}
