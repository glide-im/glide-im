package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/service/route"
	"go_im/service/rpc"
	"testing"
	"time"
)

type mockConn struct {
}

func (m *mockConn) Write(message conn.Serializable) error {
	by, _ := message.Serialize()
	fmt.Println("Write=>" + string(by))
	return nil
}

func (m *mockConn) Read(message conn.Serializable) error {
	time.Sleep(time.Hour)
	return nil
}

func (m *mockConn) Close() error { return nil }

type mockManager struct {
	manager
}

func (receiver *mockManager) addMockConn() {
	c1 := receiver.ClientConnected(&mockConn{})
	c2 := receiver.ClientConnected(&mockConn{})
	fmt.Printf("uid=%d, %d", c1, c2)
}

func TestNewClientServer(t *testing.T) {

	server := NewServer(&rpc.ServerOptions{
		Name:        "client",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8081,
		EtcdServers: etcd,
	})

	m := &mockManager{}
	client.Manager = m

	go func() {
		time.Sleep(time.Second * 3)
		m.addMockConn()
	}()
	err := server.Run()
	assert.Nil(t, err)
}

func TestRegisterService(t *testing.T) {
	err := route.RegisterService("client", etcd)
	assert.Nil(t, err)
}

func TestManager_ClientConnected(t *testing.T) {

}
