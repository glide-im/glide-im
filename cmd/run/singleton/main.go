package main

import (
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

	server.SetConnHandler(func(conn conn.Connection) {
		client.Manager.ClientConnected(conn)
	})

	group.Manager.(*group.DefaultManager).Init()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := api.RunHttpServer("0.0.0.0", 8081)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := server.Run("0.0.0.0", 8080)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
