package im

import (
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/pkg/logger"
	"time"
)

type Type int

const (
	TCP Type = iota + 1
	WebSocket
	UDP
)

type Options struct {
	SvrType Type

	ApiImpl       api.IApiHandler
	ClientMgrImpl client.IClientManager
	GroupMgrImpl  group.IGroupManager
}

type Server struct {
	opts   Options
	server conn.Server
}

func NewServer(options Options) *Server {
	dao.Init()

	ret := &Server{
		opts: options,
	}
	api.SetHandler(options.ApiImpl)
	api.MessageHandleFunc = client.EnqueueMessage

	group.Manager = options.GroupMgrImpl
	client.Manager = options.ClientMgrImpl

	manager, ok := group.Manager.(*group.DefaultManager)
	if ok {
		manager.Init()
	} else {
		logger.W("group manager not init")
	}

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	switch options.SvrType {
	case TCP:
	case WebSocket:
		ret.server = conn.NewWsServer(op)
	case UDP:
	}
	return ret
}

func (s *Server) GetConnServer() conn.Server {
	return s.server
}

func (s *Server) Serve(host string, port int) {
	s.server.SetConnHandler(onNewConn)
	err := s.server.Run(host, port)
	if err != nil {
		panic(err)
	}
}

func onNewConn(conn conn.Connection) {
	client.Manager.ClientConnected(conn)
}
