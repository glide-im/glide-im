package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/service/api/pb"
	"go_im/service/rpc"
	"time"
)

type Client struct {
	rpc pb.ApiServiceClient
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	ret := &Client{}
	ret.BaseClient = rpc.NewBaseClient(options)
	ret.Init(options)
	api.SetImpl(ret)
	return ret
}

func (c *Client) Handle(uid int64, message *message.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	m := pb.Message{
		Seq:    message.Seq,
		Action: string(message.Action),
		Data:   message.Data,
	}
	testFunc, err := c.rpc.Handle(ctx, &pb.HandleRequest{
		Uid:     uid,
		Message: &m,
	})
	if err != nil {
		panic(err)
	}
	if testFunc.GetOk() {

	}
}

func (c *Client) Run() error {
	logger.D("gRPC Api client run")
	err := c.BaseClient.Run()
	logger.D("gRPC Api client connect to %s complete, state=%s", c.Conn.Target(), c.Conn.GetState())
	c.rpc = pb.NewApiServiceClient(c.Conn)
	return err
}
