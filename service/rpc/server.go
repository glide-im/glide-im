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

type ServerOptions struct {
	Name           string
	Network        string
	Addr           string
	Port           int
	MaxRecvMsgSize int
	MaxSendMsgSize int
	EtcdServers    []string
}

type BaseServer struct {
	Srv *server.Server

	options      *ServerOptions
	etcdRegister *serverplugin.EtcdV3RegisterPlugin
	reg          []func(srv *BaseServer) error
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

func (s *BaseServer) Register(name string, sv interface{}) {
	s.reg = append(s.reg, func(srv *BaseServer) error {
		return srv.Srv.RegisterName(name, sv, "")
	})
}

func (s *BaseServer) Run() error {

	addr := fmt.Sprintf("%s:%d", s.options.Addr, s.options.Port)
	s.etcdRegister.ServiceAddress = s.options.Network + "@" + addr

	err := s.etcdRegister.Start()
	if err != nil {
		return err
	}
	s.Srv.Plugins.Add(s.etcdRegister)

	for _, f := range s.reg {
		if er := f(s); er != nil {
			return er
		}
	}

	return s.Srv.Serve(s.options.Network, addr)
}
