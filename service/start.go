package main

import (
	"go_im/im"
	api2 "go_im/im/api"
	client2 "go_im/im/client"
	"go_im/im/conn"
	group2 "go_im/im/group"
	"go_im/pkg/logger"
	"go_im/service/api"
	"go_im/service/client"
	"go_im/service/group"
	"go_im/service/rpc"
	"math"
	"os"
	"strconv"
	"sync"
)

type ServerType int

const (
	_ ServerType = iota
	TypeApiService
	TypeClientService
	TypeGroupService
)

const (
	PortClientSrv = 5555
	PortApiSrv    = 5556
	PortGroupSrv  = 5557
)

var defaultSrvOpts = rpc.ServerOptions{
	Network:        "tcp",
	Addr:           "localhost",
	MaxRecvMsgSize: math.MaxInt32,
	MaxSendMsgSize: math.MaxInt32,
}

var defaultCliOpts = rpc.ClientOptions{
	Addr: "localhost",
}

var wait = new(sync.WaitGroup)

func main() {

	t, _ := strconv.Atoi(os.Args[1])
	var sType = ServerType(t)
	go runApiService(sType)
	go runClientService(sType)
	go runGroupService(sType)
	wait.Wait()
}

func run(srv rpc.Runnable) {
	wait.Add(1)
	err := srv.Run()
	if err != nil {
		logger.E("grpc run err", err)
	}
	wait.Done()
}

func runApiService(t ServerType) {
	if TypeApiService == t {
		options := defaultSrvOpts
		options.Port = PortApiSrv
		server := api.NewServer(&options)
		api2.SetImpl(im.NewApiRouter())
		run(server)
	} else {
		clientOpts := defaultCliOpts
		clientOpts.Port = PortApiSrv
		c := api.NewClient(&clientOpts)
		api2.SetImpl(c)
		run(c)
	}
}

func runClientService(t ServerType) {
	if TypeClientService == t {
		options := defaultSrvOpts
		options.Port = PortClientSrv
		server := client.NewServer(&options)
		mgr := im.NewClientManager()
		wsServer := conn.NewWsServer(nil)
		wsServer.SetConnHandler(func(conn conn.Connection) {
			mgr.ClientConnected(conn)
		})
		client2.Manager = mgr
		go func() {
			err := wsServer.Run("localhost", 8080)
			if err != nil {
				logger.E("start ws server err", err)
			}
		}()
		run(server)
	} else {
		clientOpts := defaultCliOpts
		clientOpts.Port = PortClientSrv
		c := client.NewClient(&clientOpts)
		client2.Manager = c
		run(c)
	}
}

func runGroupService(t ServerType) {
	if TypeGroupService == t {
		options := defaultSrvOpts
		options.Port = PortGroupSrv
		server := group.NewServer(&options)
		group2.Manager = im.NewGroupManager()
		run(server)
	} else {
		clientOpts := defaultCliOpts
		clientOpts.Port = PortGroupSrv
		c := group.NewClient(&clientOpts)
		group2.Manager = c
		run(c)
	}
}
