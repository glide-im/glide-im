package rpc

import (
	"context"
	"fmt"
	etcd_cli "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type Cli interface {
	Call(ctx context.Context, fn string, request, reply interface{}) error
	Broadcast(fn string, request, reply interface{}) error
	Run() error
	Close() error
}

type ClientOptions struct {
	client.Option

	Addr        string
	Port        int
	Name        string
	EtcdServers []string
	Selector    client.Selector
}

type BaseClient struct {
	cli     client.XClient
	options *ClientOptions
	id      string
}

func NewBaseClient(options *ClientOptions) (*BaseClient, error) {
	ret := &BaseClient{
		options: options,
	}

	var discovery client.ServiceDiscovery
	var err error

	if options.EtcdServers != nil {
		discovery, err = etcd_cli.NewEtcdV3Discovery(BaseServicePath, options.Name, options.EtcdServers, false, nil)
		if err != nil {
			return nil, err
		}
	} else {
		srv := fmt.Sprintf("%s@%s:%d", "tcp", options.Addr, options.Port)
		discovery, _ = client.NewPeer2PeerDiscovery(srv, "")
	}

	if options.SerializeType == protocol.SerializeNone {
		// using protobuffer serializer by default
		options.SerializeType = protocol.ProtoBuffer
	}
	ret.cli = client.NewXClient(options.Name, client.Failtry, client.RoundRobin, discovery, options.Option)

	if options.Selector != nil {
		ret.cli.SetSelector(options.Selector)
	} else {
		// using round robbin selector by default
		ret.cli.SetSelector(NewRoundRobinSelector())
	}
	return ret, nil
}

func (c *BaseClient) Call2(fn string, arg interface{}, reply interface{}) error {
	return c.Call(context.Background(), fn, arg, reply)
}

func (c *BaseClient) Broadcast(fn string, request, reply interface{}) error {
	return c.cli.Broadcast(context.Background(), fn, request, reply)
}

func (c *BaseClient) Call(ctx context.Context, fn string, arg interface{}, reply interface{}) error {
	err := c.cli.Call(ctx, fn, arg, reply)
	return err
}

func (c *BaseClient) Run() error {
	return nil
}

func (c *BaseClient) Close() error {
	return c.cli.Close()
}
