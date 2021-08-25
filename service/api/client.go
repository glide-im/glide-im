package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/api/rpc"
	rpc2 "go_im/service/rpc"
	"time"
)

type Client struct {
	rpc rpc.ApiServiceClient
	*rpc2.BaseClient
}

func NewClient(options *rpc2.ClientOptions) *Client {
	ret := &Client{}
	ret.BaseClient = rpc2.NewBaseClient(options)
	ret.Init(options)
	api.SetImpl(ret)
	return ret
}

func (c *Client) Handle(uid int64, message *message.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	m := rpc.Message{
		Seq:    message.Seq,
		Action: string(message.Action),
		Data:   message.Data,
	}
	testFunc, err := c.rpc.Handle(ctx, &rpc.Request{
		Uid:     uid,
		Message: &m,
	})
	if err != nil {
		panic(err)
	}
	if testFunc.Ok {

	}
}

func (c *Client) Run() error {
	err := c.Connect()
	c.rpc = rpc.NewApiServiceClient(c.Conn)
	return err
}
