package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_im/im/client"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/service/route"
	"go_im/service/rpc"
	"testing"
	"time"
)

type mockConn struct {
}

func (m *mockConn) Write(data []byte) error {
	fmt.Println("MockConnWrite=>" + string(data))
	return nil
}

func (m *mockConn) Read() ([]byte, error) {
	select {}
}

func (m *mockConn) Close() error { return nil }

type mockManager struct {
	*manager
}

func newMockManager(opts *rpc.ServerOptions) *mockManager {
	n, _ := newManager(etcd, fmt.Sprintf("%s@%s:%d", opts.Network, opts.Addr, opts.Port))
	return &mockManager{manager: n}
}

func (receiver *mockManager) addMockConn() []int64 {
	return []int64{
		receiver.ClientConnected(&mockConn{}),
		receiver.ClientConnected(&mockConn{}),
	}
}

func TestRegisterService(t *testing.T) {
	err := route.RegisterService("client", etcd)
	assert.Nil(t, err)
}

func TestServer_EnqueueMessage(t *testing.T) {
	opts := &rpc.ClientOptions{
		Name:        ServiceName,
		EtcdServers: etcd,
	}
	cli, err := NewClientByRouter(opts)
	assert.Nil(t, err)
	cli.EnqueueMessage(1000000301026, 0, message.NewMessage(1, "api.app.echo", "hello world"))
}

func TestNewServer(t *testing.T) {
	type args struct {
		options *rpc.ServerOptions
	}
	getArgs := func(port int) args {
		return args{options: &rpc.ServerOptions{
			Name:        "client",
			Network:     "tcp",
			Addr:        "127.0.0.1",
			Port:        port,
			EtcdServers: etcd,
		}}
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{name: "WithConn_1", args: getArgs(8081)},
		{name: "WithConn_2", args: getArgs(8082)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid.Mock()
			m := newMockManager(tt.args.options)
			s := NewServer(tt.args.options)
			client.Manager = m
			go func() {
				time.Sleep(time.Second * 1)
				t.Logf("conn id for test=%v", m.addMockConn())
			}()
			assert.Nil(t, s.Run())
		})
	}
}
