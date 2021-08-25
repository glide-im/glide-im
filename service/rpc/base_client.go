package rpc

import (
	"fmt"
	"google.golang.org/grpc"
)

type ClientOptions struct {
	Addr string
	Port int
}

type BaseClient struct {
	Conn *grpc.ClientConn

	options *ClientOptions
}

func NewBaseClient(options *ClientOptions) *BaseClient {
	ret := &BaseClient{}
	ret.Init(options)
	return ret
}

func (b *BaseClient) Init(options *ClientOptions) {
	if options == nil {
		b.options = &ClientOptions{
			Addr: "localhost",
			Port: 5555,
		}
	} else {
		b.options = options
	}
}

func (b *BaseClient) Connect() error {
	if b.options == nil {
		b.Init(nil)
	}
	var err error
	target := fmt.Sprintf("%s:%d", b.options.Addr, b.options.Port)
	b.Conn, err = grpc.Dial(target, grpc.WithInsecure())
	return err
}
