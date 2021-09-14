package rpc

import (
	"context"
	"fmt"
	client3 "github.com/rpcxio/rpcx-etcd/client"
	client2 "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type ClientOptions struct {
	client2.Option

	Addr        string
	Port        int
	Name        string
	EtcdServers []string
}

type BaseClient struct {
	cli     client2.XClient
	options *ClientOptions
	id      string
}

func NewBaseClient(options *ClientOptions) *BaseClient {
	ret := &BaseClient{
		options: options,
		id:      fmt.Sprintf("%s@%s:%d", "", "", 1),
	}
	etcd, err := client3.NewEtcdV3Discovery(BaseServicePath, options.Name, options.EtcdServers, false, nil)
	if err != nil {

	}
	if options.SerializeType == protocol.SerializeNone {
		options.SerializeType = protocol.ProtoBuffer
	}
	ret.cli = client2.NewXClient(options.Name, client2.Failtry, client2.RoundRobin, etcd, options.Option)
	return ret
}

func (c *BaseClient) Call(fn string, arg interface{}, reply interface{}) error {
	//ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{"call_from_client_server": c.id})
	//ctx = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
	return c.Call2(context.Background(), fn, arg, reply)
}

func (c *BaseClient) Call2(ctx context.Context, fn string, arg interface{}, reply interface{}) error {
	err := c.cli.Call(ctx, fn, arg, reply)
	return err
}

func (c *BaseClient) Run() error {
	c.cli.SetSelector(NewHostRouter())
	return nil
}

func (c *BaseClient) Close() error {
	return c.cli.Close()
}
