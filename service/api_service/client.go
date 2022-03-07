package api_service

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/protobuff/gen/pb_rpc"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
	api.Handler = ret
	return ret, nil
}

func (c *Client) Echo(uid int64, message *message.Message) *pb_rpc.Response {
	//m := pb.Message{
	//	Seq:    message.Seq,
	//	Action: string(message.Action),
	//}
	arg := &pb_rpc.ApiHandleRequest{
		Uid: uid,
		//Message: &m.,
	}

	resp := &pb_rpc.Response{
		Ok:      false,
		Message: "",
	}
	err := c.Call(rpc.NewCtx(), "Echo", arg, resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *Client) Handle(uid int64, device int64, message *message.Message) error {
	ctx := context.WithValue(context.Background(), "from_gate", "node_id")

	request := pb_rpc.ApiHandleRequest{
		Uid:     uid,
		Device:  device,
		Message: message.CommMessage,
	}
	err := c.Call(ctx, "Handle", &request, &emptypb.Empty{})
	return err
}
