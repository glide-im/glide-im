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

	manager := group.NewDefaultManager()
	group.SetInterfaceImpl(manager)
	manager.Init()

	errCh := make(chan error)

	go func() {
		errCh <- server.Run(config.IMRpcServer.Addr, config.IMRpcServer.Port)
	}()

	go func() {
		options := rpc.ServerOptions{}

		srvName := config.IMRpcServer.Name
		etcd := config.IMRpcServer.Etcd

		if srvName != "" && len(etcd) > 0 {
			options.Name = srvName
			options.EtcdServers = etcd
			logger.D("start im rpc server by etcd")
		} else {

			addr := config.IMRpcServer.Addr
			port := config.IMRpcServer.Port
			if addr == "" || port == 0 {
				errCh <- errors.New("rpc server addr or port is empty")
				return
			}
			options.Addr = addr
			options.Port = port
			logger.D("start im rpc server by addr")
		}

		rpcServer := im_service.NewServer(&options)
		errCh <- rpcServer.Run()
	}()

	panic(<-errCh)
}
