package service

import (
	"github.com/BurntSushi/toml"
	"go_im/pkg/rpc"
	"sync"
)

var c *Configs
var loadErr error
var once = sync.Once{}

func GetConfig() (*Configs, error) {
	once.Do(func() {
		c = &Configs{}
		_, loadErr = toml.DecodeFile("example_config.toml", c)
	})
	return c, loadErr
}

type ServerConfig struct {
	Addr    string
	Port    int
	Network string
	Name    string
	SrvID   string
}

func (s *ServerConfig) ToServerOptions(etcd []string) *rpc.ServerOptions {
	return &rpc.ServerOptions{
		Name:        s.Name,
		Network:     s.Network,
		Addr:        s.Addr,
		Port:        s.Port,
		EtcdServers: etcd,
	}
}

type ClientConfig struct {
	Retries           int32
	IdleTimeout       int64
	ConnectTimeout    int64
	Heartbeat         bool
	HeartbeatInterval int64

	Name        string
	EtcdServers []string
	// optional when use service discovery
	Addr string
	Port int32
}

func (c *ClientConfig) ToClientOptions() *rpc.ClientOptions {
	return &rpc.ClientOptions{
		Addr:        c.Addr,
		Port:        int(c.Port),
		Name:        c.Name,
		EtcdServers: c.EtcdServers,
	}
}

type ApiConfig struct {
	Server *ServerConfig
	Client *ClientConfig
}

type GatewayConfig struct {
	Server *ServerConfig
	Client *ClientConfig
}

type GroupMessagingConfig struct {
	Server *ServerConfig
	Client *ClientConfig
}

type DispatchConfig struct {
	Server *ServerConfig
	Client *ClientConfig
}

type MessageRouterConfig struct {
	Server *ServerConfig
	Client *ClientConfig
}

type EtcdConfig struct {
	Servers []string
}

type NsqConfig struct {
	Lookup string
	Nsqd   string
}

type Configs struct {
	Etcd *EtcdConfig
	Nsq  *NsqConfig

	Dispatch       *DispatchConfig
	Api            *ApiConfig
	Gateway        *GatewayConfig
	GroupMessaging *GroupMessagingConfig
	MessageRouter  *MessageRouterConfig
}
