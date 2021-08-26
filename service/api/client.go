package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
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
	err := c.Connect()
	c.rpc = pb.NewApiServiceClient(c.Conn)
	return err
}
