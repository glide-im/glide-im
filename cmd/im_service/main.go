package main

import (
	"errors"
	"go_im/config"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/messaging"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"go_im/service/im_service"
	"time"
)

func main() {

	err := config.Load()
	if err != nil {
		panic(err)
	}
	db.Init()
	dao.Init()

	var server conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	client.SetMessageHandler(messaging.HandleMessage)
	server = conn.NewWsServer(op)

	cm := client.NewDefaultManager()
	server.SetConnHandler(func(conn conn.Connection) {
		cm.ClientConnected(conn)
	})

	client.SetInterfaceImpl(cm)

	if config.IMRpcServer.EnableGroup {
		manager := group.NewDefaultManager()
		group.SetInterfaceImpl(manager)
		manager.Init()
	}

	errCh := make(chan error)

	go func() {
		errCh <- server.Run(config.WsServer.Addr, config.WsServer.Port)
	}()

	go func() {
		options := rpc.ServerOptions{}

		options.Addr = config.IMRpcServer.Addr
		options.Port = config.IMRpcServer.Port
		options.Name = config.IMRpcServer.Name
		options.Network = config.IMRpcServer.Network
		options.EtcdServers = config.IMRpcServer.Etcd

		if options.Name != "" && len(options.EtcdServers) > 0 {
			logger.D("start im rpc server by etcd")
		} else {
			if options.Addr == "" || options.Port == 0 {
				errCh <- errors.New("rpc server addr or port is empty")
				return
			}
			logger.D("start im rpc server by addr")
		}
		err2 := im_service.RunServer(&options)
		errCh <- err2
	}()

	panic(<-errCh)
}
