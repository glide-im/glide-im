package api_service

import (
	"context"
	"go_im/im/message"
	rpc2 "go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	rpc2.Cli
}

func NewClient(options *rpc2.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = rpc2.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
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
	err := c.Call(rpc2.NewCtx(), "Echo", arg, resp)
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
		Message: message.GetProtobuf(),
	}
	err := c.Call(ctx, "Handle", &request, &emptypb.Empty{})
	return err
}
