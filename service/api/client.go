package api

import (
	"context"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/pb"
	"go_im/service/route"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/anypb"
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

func NewClientByRouter(srvId string, rtOpts *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error
	ret.Cli, err = route.NewRouter(srvId, rtOpts)
	if err != nil {
		return nil, err
	}
	api.Handler = ret
	return ret, nil
}

func (c *Client) Echo(uid int64, message *message.Message) *pb.Response {
	//m := pb.Message{
	//	Seq:    message.Seq,
	//	Action: string(message.Action),
	//}
	arg := &pb.HandleRequest{
		Uid: uid,
		//Message: &m.,
	}

	resp := &pb.Response{
		Ok:      false,
		Message: "",
	}
	err := c.Call(rpc.NewCtx(), "Echo", arg, resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *Client) Handle(uid int64, device int64, message *message.Message) {
	ctx := context.WithValue(context.Background(), "from_gate", "node_id")

	any, err2 := anypb.New(message)
	if err2 != nil {
		return
	}
	request := pb.HandleRequest{
		Uid:     uid,
		Device:  device,
		Message: any,
	}
	err := c.Call(ctx, "Handle", &request, &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
}
