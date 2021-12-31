package main

import (
	"go_im/im/api"
	"go_im/im/api/apidep"
	"go_im/im/dao"
	"go_im/pkg/db"
)

func main() {

	db.Init()
	dao.Init()

	apidep.GroupManager = &apidep.MockGroupManager{}
	apidep.ClientManager = &apidep.MockClientManager{}
	err := api.RunHttpServer("0.0.0.0", 8081)

	if err != nil {
		panic(err)
	}
}
