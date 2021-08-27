package rpc

import (
	"context"
	"errors"
	"fmt"
	"go_im/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientOptions struct {
	Addr string
	Port int
}

type BaseClient struct {
	Conn *grpc.ClientConn

	AppId   int64
	Options *ClientOptions
}

func NewBaseClient(options *ClientOptions) *BaseClient {
	ret := &BaseClient{}
	ret.Init(options)
	return ret
}

func (b *BaseClient) Init(options *ClientOptions) {
	if options == nil {
		b.Options = &ClientOptions{
			Addr: "localhost",
			Port: 5555,
		}
	} else {
		b.Options = options
	}
}

func (b *BaseClient) unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if b.Conn.GetState() != connectivity.Ready {
		return errors.New("client is not connect to the server")
	}
	logger.D("rpc client call method: %s", method)
	return nil
}

func (b *BaseClient) Run() error {
	if b.Options == nil {
		b.Init(nil)
	}
	var err error
	target := fmt.Sprintf("%s:%d", b.Options.Addr, b.Options.Port)
	b.Conn, err = grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(b.unaryInterceptor))

	return err
}

type ProxyRpcClientConn struct {
	consumer grpc.ClientConnInterface
}

func (p ProxyRpcClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	panic("implement me")
}

func (p ProxyRpcClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	panic("implement me")
}
