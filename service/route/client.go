package route

import (
	"context"
	"errors"
	"github.com/smallnest/rpcx/share"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type Client struct {
	rpc.Cli
}

func NewClient(options *rpc.ClientOptions) (*Client, error) {
	c, err := rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		Cli: c,
	}
	return cli, nil
}

func (c *Client) Unregister(srvId string) error {
	return c.Broadcast("Unregister", &pb.UnRegisterReq{SrvId: srvId}, &emptypb.Empty{})
}

func (c *Client) Register(param *pb.RegisterRtReq, reply *emptypb.Empty) error {
	return c.Broadcast("Register", param, reply)
}

func (c *Client) SetTag(srvId, tag, value string) error {
	req := &pb.SetTagReq{
		Tag:   tag,
		SrvId: srvId,
		Value: value,
	}
	return c.Broadcast("SetTag", req, &emptypb.Empty{})
}

func (c *Client) RemoveTag(srvId, tag string) error {
	return c.Broadcast("RemoveTag", &pb.ClearTagReq{
		SrvId: srvId,
		Tag:   tag,
	}, &emptypb.Empty{})
}

func (c *Client) GetAllTag(srvId string) (map[string]string, error) {
	resp := &pb.AllTagResp{Tags: map[string]string{}}
	err := c.Call(context.Background(), "GetAllTag", &pb.AllTagReq{SrvId: srvId}, resp)
	return resp.GetTags(), err
}

func (c *Client) Route(ctx context.Context, extra map[string]string, target string, request, reply interface{}) error {

	split := strings.Split(target, ".")
	if len(split) != 2 {
		return errors.New("参数 target 格式错误, (srvId.func).() 例子: api.Handle")
	}

	var reqParam *anypb.Any
	var err error

	if p, ok := request.(proto.Message); ok {
		reqParam, err = anypb.New(p)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request must be proto.Message")
	}

	routeReq := &pb.RouteReq{
		SrvId:  split[0],
		Fn:     split[1],
		Params: reqParam,
		Extra:  extra,
	}
	routeReply := &pb.RouteReply{}
	err = c.Call(ctx, "Route", routeReq, routeReply)

	if err != nil {
		return err
	}
	if resp, ok := reply.(proto.Message); ok {
		if !routeReply.GetReply().MessageIs(resp) {
			return errors.New("route reply message not matched to source reply")
		}
		err = routeReply.GetReply().UnmarshalTo(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) RouteByTag(target, tag string, request, reply interface{}) error {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{ExtraTag: tag})
	return c.Route(ctx, map[string]string{}, target, request, reply)
}

func RegisterService(srvId string, etcd []string) error {
	cli, err := NewClient(&rpc.ClientOptions{
		Name:        ServiceName,
		EtcdServers: etcd,
	})
	defer func() {
		_ = cli.Close()
	}()
	if err != nil {
		return err
	}

	req := &pb.RegisterRtReq{
		SrvId:           srvId,
		DiscoverySrvUrl: etcd,
	}
	return cli.Register(req, &emptypb.Empty{})
}

/////////////////////////////////////////////////////////////////////////////////////////

type RouterCli struct {
	rt    *Client
	srvId string
}

func NewRouter(srvId string, routeOpts *rpc.ClientOptions) (*RouterCli, error) {
	routeOpts.Name = ServiceName
	c, err := NewClient(routeOpts)
	if err != nil {
		return nil, err
	}
	return &RouterCli{
		rt:    c,
		srvId: srvId,
	}, nil
}

func (r *RouterCli) Call(ctx context.Context, fn string, request, reply interface{}) error {
	path := r.srvId + "." + fn
	return r.rt.Route(ctx, map[string]string{}, path, request, reply)
}

func (r *RouterCli) Broadcast(fn string, request, reply interface{}) error {
	return errors.New("broadcast on proxy mode is unsupported")
}

func (r *RouterCli) Run() error {
	return r.rt.Run()
}

func (r *RouterCli) Close() error {
	return r.rt.Close()
}
