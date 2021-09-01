package api

import (
	"context"
	"github.com/rcrowley/go-metrics"
	client3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	client2 "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"go_im/im"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/group"
	"go_im/service/rpc"
	"math"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {

	api.SetImpl(im.NewApiRouter())
	client.Manager = im.NewClientManager()
	group.Manager = im.NewGroupManager()

	op := rpc.ServerOptions{
		Network:        "tcp",
		Addr:           "localhost",
		Port:           5555,
		MaxRecvMsgSize: math.MaxInt32,
		MaxSendMsgSize: math.MaxInt32,
	}
	server := NewServer(&op)
	err := server.Run()
	panic(err)
}

type Ts struct {
}

type TestArg struct {
	Name string
}

type TestReply struct {
	Name string
}

func (a *Ts) Handle2(ctx context.Context, param *TestArg, reply *TestReply) error {
	reply.Name = param.Name + "==="
	return nil
}

func TestServer_Handle2(t *testing.T) {
	s := server.NewServer()

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@127.0.0.1:8972",
		EtcdServers:    []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
		BasePath:       "/test",
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	_ = r.Start()
	s.Plugins.Add(r)

	//_ = s.RegisterName("api", new(Server), "")
	_ = s.RegisterName("ts", new(Ts), "")
	_ = s.Serve("tcp", "127.0.0.1:8972")
}

func TestServer_Handle1(t *testing.T) {
	s := server.NewServer()

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@127.0.0.1:8973",
		EtcdServers:    []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
		BasePath:       "/test",
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	_ = r.Start()
	s.Plugins.Add(r)

	//_ = s.RegisterName("api", new(Server), "")
	_ = s.RegisterName("ts", new(Ts), "")
	_ = s.Serve("tcp", "127.0.0.1:8973")
}

func Test3(t *testing.T) {

	etcdr, _ := client3.NewEtcdV3Discovery("/test", "ts",
		[]string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"},
		false, nil)

	//discovery, _ := client2.NewPeer2PeerDiscovery("tcp@127.0.0.1:8972", "")
	xClient := client2.NewXClient("ts", client2.Failtry, client2.RoundRobin, etcdr, client2.DefaultOption)
	defer xClient.Close()
	xClient.SetSelector(rpc.NewHostRouter())
	r := TestReply{}
	err2 := xClient.Call(context.Background(), "Handle2", &TestArg{Name: "hello"}, &r)
	if err2 != nil {
		t.Log(err2)
	}
	t.Log(r.Name)
}
