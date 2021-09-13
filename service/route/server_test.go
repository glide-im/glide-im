package route

import (
	"context"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

var etcdSrv = []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"}

func TestNewServer(t *testing.T) {

	op := rpc.ServerOptions{
		Name:        "route",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8975,
		EtcdServers: etcdSrv,
	}

	routeServer := NewServer(&op)
	err := routeServer.Run()
	t.Error(err)
}

func TestClient_Register(t *testing.T) {

	cli := newClient()

	err := cli.Register(&pb.RegisterRtReq{
		SrvId:           "api",
		SrvName:         "api",
		RoutePolicy:     1,
		DiscoverySrvUrl: etcdSrv,
		DiscoveryType:   1,
	}, &emptypb.Empty{})

	if err != nil {
		t.Error(err)
	}
}

func TestClient_Route(t *testing.T) {
	cli := newClient()
	err := cli.Route(context.Background(), &pb.RouteReq{
		SrvId: "api",
		Fn:    "Handle",
	}, &pb.Response{})
	if err != nil {
		t.Error(err)
	}
}

func newClient() *Client {
	client := NewClient(&rpc.ClientOptions{
		Name:        "route",
		EtcdServers: etcdSrv,
	})
	err := client.Run()
	if err != nil {
		panic(err)
	}
	return client
}
