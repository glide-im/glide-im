package rpc

import (
	"context"
	client3 "github.com/rpcxio/rpcx-etcd/client"
	client2 "github.com/smallnest/rpcx/client"
)

type BaseClient struct {
	cli     client2.XClient
	options *ClientOptions
}

func NewBaseClient(options *ClientOptions) *BaseClient {
	ret := &BaseClient{options: options}
	etcd, err := client3.NewEtcdV3Discovery(BaseServicePath, options.Name, options.EtcdServers, false, nil)
	if err != nil {

	}
	ret.cli = client2.NewXClient(options.Name, client2.Failtry, client2.RoundRobin, etcd, client2.DefaultOption)
	return ret
}

func (c *BaseClient) Call(fn string, arg interface{}, reply interface{}) error {
	return c.cli.Call(context.Background(), fn, arg, reply)
}

func (c *BaseClient) Run() error {
	c.cli.SetSelector(NewHostRouter())
	return nil
}

func (c *BaseClient) Close() error {
	return c.cli.Close()
}
