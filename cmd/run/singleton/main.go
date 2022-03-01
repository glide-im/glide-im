package main

import (
	"go_im/config"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/messaging"
	"go_im/pkg/db"
	"sync"
	"time"
)

func main() {
	Run()
}

func Run() {
	db.Init()
	dao.Init()

	var server conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	messaging.Init()
	server = conn.NewWsServer(op)

	cm := client.NewDefaultManager()
	server.SetConnHandler(func(conn conn.Connection) {
		cm.ClientConnected(conn)
	})

	client.SetInterfaceImpl(cm)

	manager := group.NewDefaultManager()
	group.SetInterfaceImpl(manager)
	manager.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := api.RunHttpServer(config.ApiHttpService.Addr, config.ApiHttpService.Port)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := server.Run(config.IMService.Addr, config.IMService.Port)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
