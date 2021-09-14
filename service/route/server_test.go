package route

import (
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
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
	go func() {
		time.Sleep(time.Second * 1)
		TestClient_Register(t)
	}()
	err := routeServer.Run()
	t.Error(err)
}

func TestClient_Register(t *testing.T) {

	cli := newClient()
	defer cli.Close()
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
	defer cli.Close()
	err := cli.Invoke("api.Handle", &pb.HandleRequest{
		Uid:     1,
		Message: nil,
	}, &emptypb.Empty{})
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
