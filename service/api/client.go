package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/route"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) *Client {
	ret := &Client{}
	ret.Cli = rpc.NewBaseClient(options)
	api.SetImpl(ret)
	return ret
}

func NewClientByRouter(srvId string, rtOpts *rpc.ClientOptions) *Client {
	ret := &Client{}
	ret.Cli = route.NewRouter(srvId, rtOpts)
	api.SetImpl(ret)
	return ret
}

func (c *Client) Handle(uid int64, message *message.Message) {
	m := pb.Message{
		Seq:    message.Seq,
		Action: string(message.Action),
		Data:   message.Data,
	}
	arg := &pb.HandleRequest{
		Uid:     uid,
		Message: &m,
	}

	err := c.Call(context.Background(), "Handle", arg, &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
}
