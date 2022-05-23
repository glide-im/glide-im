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

	Options      *ServerOptions
	etcdRegister *serverplugin.EtcdV3RegisterPlugin
	reg          []func(srv *BaseServer) error
	id           string
}

func NewBaseServer(options *ServerOptions) *BaseServer {
	ret := &BaseServer{
		Srv: server.NewServer(),
		id:  fmt.Sprintf("%s@%s:%d", options.Name, options.Addr, options.Port),
	}

	if options.Network == "" {
		options.Network = "tcp"
	}

	ret.Options = options
	if len(options.EtcdServers) != 0 {
		ret.etcdRegister = &serverplugin.EtcdV3RegisterPlugin{
			EtcdServers:    options.EtcdServers,
			BasePath:       BaseServicePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
	}
	return ret
}

func (s *BaseServer) GetServerID() string {
	if len(s.id) == 0 {
		s.id = fmt.Sprintf("%s@%s:%d", s.Options.Name, s.Options.Addr, s.Options.Port)
	}
	return s.id
}

func (s *BaseServer) Register(name string, sv interface{}) {
	s.reg = append(s.reg, func(srv *BaseServer) error {
		return srv.Srv.RegisterName(name, sv, "")
	})
}

func (s *BaseServer) Run() error {

	addr := fmt.Sprintf("%s:%d", s.Options.Addr, s.Options.Port)

	if s.etcdRegister != nil {
		s.etcdRegister.ServiceAddress = s.Options.Network + "@" + addr

		err := s.etcdRegister.Start()
		if err != nil {
			return err
		}
		s.Srv.Plugins.Add(s.etcdRegister)
	}

	for _, f := range s.reg {
		if er := f(s); er != nil {
			return er
		}
	}

	return s.Srv.Serve(s.Options.Network, addr)
}
