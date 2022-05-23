package gateway

import (
	"fmt"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/messaging"
	"go_im/pkg/logger"
	"go_im/service"
	"go_im/service/dispatch"
	"time"
)

var gatewayRoute = "gateway"

// SetupClient init nsq client if you want to use client interface to send message
func SetupClient(configs *service.ClientConfig) error {

	options := configs.ToClientOptions()
	cli, err := NewClient(options)
	if err != nil {
		return err
	}
	client.SetInterfaceImpl(cli)
	return nil
}

// RunServer TODO 2022-3-24 run nsq
func RunServer(configs *service.Configs) error {

	dispatchCliOpts := configs.Dispatch.Client.ToClientOptions()
	dispatchCliOpts.EtcdServers = configs.Etcd.Servers

	dispatchService, err := dispatch.NewClient(dispatchCliOpts)
	if err != nil {
		return err
	}

	gatewaySrvConfig := configs.Gateway.Server

	gatewayRoute = fmt.Sprintf("%s@%s:%d", gatewaySrvConfig.Network, gatewaySrvConfig.Addr, gatewaySrvConfig.Port)

	var imServer conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	client.SetMessageHandler(messaging.HandleMessage)
	imServer = conn.NewWsServer(op)

	cm := client.NewDefaultManager()
	imServer.SetConnHandler(func(conn conn.Connection) {
		id := cm.ClientConnected(conn)
		e := dispatchService.UpdateGatewayRoute(id, gatewayRoute)
		if e != nil {
			logger.E("update gateway route error: %v", e)
		}
	})

	client.SetInterfaceImpl(cm)

	ch := make(chan error)

	go func() {
		logger.D("im server starting")
		ch <- imServer.Run("0.0.0.0", 8080)
	}()

	// TODO 2022-3-24 Remove, no one use Gateway service
	go func() {
		logger.D("gateway starting")
		options := gatewaySrvConfig.ToServerOptions(configs.Etcd.Servers)
		server := NewServer(options)
		ch <- server.Run()
	}()

	return <-ch
}
