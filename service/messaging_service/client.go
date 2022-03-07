package messaging_service

import (
	"context"
	"go_im/im/message"
	"go_im/protobuff/gen/pb_rpc"
	"go_im/service/rpc"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) HandleMessage(from int64, device int64, message *message.Message) error {
	request := pb_rpc.MessagingHandleRequest{
		Id:      from,
		Device:  device,
		Message: message.CommMessage,
	}

	return c.Call(context.Background(), "UpdateMember", &request, &pb_rpc.Response{})
}
