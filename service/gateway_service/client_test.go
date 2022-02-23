package gateway_service

import (
	"github.com/stretchr/testify/assert"
	"go_im/service/rpc"
	"testing"
)

var etcd = []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"}

func TestName(t *testing.T) {
	cli, err := NewClientByRouter(&rpc.ClientOptions{
		Name:        "client",
		EtcdServers: etcd,
	})
	assert.Nil(t, err)

	cli.Run()
}
