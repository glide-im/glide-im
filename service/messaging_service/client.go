package messaging_service

import (
	"context"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
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
		Message: message.GetProtobuf(),
	}

	return c.Call(context.TODO(), "HandleMessage", &request, &pb_rpc.Response{})
}
