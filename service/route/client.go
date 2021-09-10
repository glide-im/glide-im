package route

import (
	"context"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/rpc"
)

type Client struct {
	*rpc.BaseClient
}

func NewClient(options *rpc.ClientOptions) *Client {
	return &Client{
		BaseClient: rpc.NewBaseClient(options),
	}
}

func (c *Client) Route(ctx context.Context, param *pb.RouteReq, reply *pb.Response) error {

	return nil
}

func (c *Client) RouteUserMessage(uid int64, message *message.Message) {

}

func (c *Client) RouteGroupMessage(gid int64, message *message.Message) {

}

func (c Client) GroupOnline(gid int64, rt string) {

}

func (c *Client) GroupOffline(gid int64) {

}

func (c Client) UserOnline(uid int64, message *message.Message) {

}

func (c *Client) UserOffline(uid int64) {

}
