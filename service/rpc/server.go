package rpc

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

const (
	BaseServicePath = "/im_service"
)

type BaseServer struct {
	Srv *server.Server

	serviceName  string
	options      *ServerOptions
	etcdRegister *serverplugin.EtcdV3RegisterPlugin
}

func NewBaseServer(options *ServerOptions) *BaseServer {
	ret := &BaseServer{
		Srv: server.NewServer(),
	}
	ret.options = options
	ret.etcdRegister = &serverplugin.EtcdV3RegisterPlugin{
		EtcdServers:    options.EtcdServers,
		BasePath:       BaseServicePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	return ret
}

func (s *BaseServer) Register(sv interface{}) error {
	return s.Srv.RegisterName(s.serviceName, sv, "")
}

func (s *BaseServer) Run() error {

	addr := fmt.Sprintf("%s:%d", s.options.Addr, s.options.Port)
	s.etcdRegister.ServiceAddress = s.options.Network + "@" + addr

	err := s.etcdRegister.Start()
	if err != nil {
		return err
	}
	s.Srv.Plugins.Add(s.etcdRegister)
	return s.Srv.Serve(s.options.Network, addr)
}
