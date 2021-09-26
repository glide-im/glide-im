package route

import (
	"github.com/stretchr/testify/assert"
	"go_im/service/rpc"
	"testing"
)

var etcdSrv = []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"}

func TestNewRouteServer(t *testing.T) {
	type args struct{ options *rpc.ServerOptions }
	getArgs := func(port int) args {
		return args{&rpc.ServerOptions{
			Name:        ServiceName,
			Network:     "tcp",
			Addr:        "127.0.0.1",
			Port:        port,
			EtcdServers: etcdSrv,
		}}
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{name: "route_8976", args: getArgs(8976)},
		{name: "route_8977", args: getArgs(8977)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, NewServer(tt.args.options).Run())
		})
	}
}

func TestServer_GetAllTag(t *testing.T) {
	tag, err := newClient().GetAllTag("client")
	assert.Nil(t, err)
	t.Log(tag)
}
