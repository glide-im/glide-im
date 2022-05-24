package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/glide-im/glideim/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/stats"
)

type BaseGRpcClient struct {
	Conn *grpc.ClientConn

	AppId   int64
	Options *ClientOptions
}

func NewBaseGRpcClient(options *ClientOptions) *BaseGRpcClient {
	ret := &BaseGRpcClient{}
	ret.Init(options)
	return ret
}

func (b *BaseGRpcClient) Init(options *ClientOptions) {
	if options == nil {
		b.Options = &ClientOptions{
			Addr: "localhost",
			Port: 5555,
		}
	} else {
		b.Options = options
	}
}

func (b *BaseGRpcClient) unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if b.Conn.GetState() != connectivity.Ready {
		return errors.New("client is not connect to the server")
	}
	logger.D("rpc client call method: %s", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		logger.E("rpc client method call error", err)
	}
	logger.D("response=%v", reply)
	return err
}

func (b *BaseGRpcClient) Run() error {

	if b.Options == nil {
		b.Init(nil)
	}
	var err error
	target := fmt.Sprintf("%s:%d", b.Options.Addr, b.Options.Port)

	b.Conn, err = grpc.Dial(target,
		grpc.WithInsecure(), // insecure connection
		grpc.WithBlock(),    // blocking until dial success
		grpc.WithUnaryInterceptor(b.unaryInterceptor),
		grpc.WithStatsHandler(newStateHandler()),
		grpc.WithUserAgent("client-id: none"))

	return err
}

type statsHandler struct {
}

func newStateHandler() *statsHandler {
	return &statsHandler{}
}

func (h *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return context.TODO()
}

func (h *statsHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {

}

func (h *statsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return context.TODO()
}

func (h *statsHandler) HandleConn(ctx context.Context, connStats stats.ConnStats) {

}
