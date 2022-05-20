package main

import (
	"go_im/im/api"
	"go_im/im/dao"
	"go_im/pkg/db"
)

func main() {

	db.Init()
	dao.Init()

	api.MockDep()
	err := api.RunHttpServer("0.0.0.0", 8081)

	if err != nil {
		panic(err)
	}
}
