package service

import (
	"github.com/BurntSushi/toml"
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

	Api            *ApiConfig
	Gateway        *GatewayConfig
	GroupMessaging *GroupMessagingConfig
	MessageRouter  *MessageRouterConfig
}
