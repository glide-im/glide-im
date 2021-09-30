package client

import (
	"fmt"
	"go_im/im"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/service/route"
	"go_im/service/rpc"
)

type manager struct {
	appId  int64
	m      *im.ClientManagerImpl
	router *route.Client
	myAddr string
}

func newManager(etcd []string, myAddr string) (*manager, error) {
	ret := &manager{}
	ret.myAddr = myAddr
	ret.m = im.NewClientManager()
	options := &rpc.ClientOptions{
		Name:        route.ServiceName,
		EtcdServers: etcd,
	}
	var err error
	ret.router, err = route.NewClient(options)
	return ret, err
}

func (m *manager) AddClient(uid int64, cs client.IClient) {
	uidTag := fmt.Sprintf("uid_%d", uid)
	err := m.router.SetTag("client", uidTag, m.myAddr)
	if err != nil {
		logger.E("set route tag error", err)
	} else {
		m.m.AddClient(uid, cs)
	}
}

func (m *manager) ClientConnected(conn conn.Connection) int64 {
	uid := m.m.ClientConnected(conn)
	tag := fmt.Sprintf("uid_%d_%d", uid, client.DeviceUnknown)
	err := m.router.SetTag("client", tag, m.myAddr)
	if err != nil {
		logger.E("set route tag error", err)
		return 0
	}
	return uid
}

func (m *manager) ClientSignIn(oldUid int64, uid int64, device int64) {
	err := m.router.RemoveTag("client", fmt.Sprintf("uid_%d_%d", oldUid, client.DeviceUnknown))
	if err != nil {

	}

	tag := fmt.Sprintf("uid_%d_%d", uid, device)
	err = m.router.SetTag("client", tag, m.myAddr)
	if err != nil {

	}
	m.m.ClientSignIn(oldUid, uid, device)
}

func (m *manager) ClientLogout(uid int64, device int64) {
	err := m.router.RemoveTag("client", fmt.Sprintf("uid_%d_%d", uid, device))
	if err != nil {

	}
	m.m.ClientLogout(uid, device)
}

func (m *manager) EnqueueMessage(uid int64, message *message.Message) {
	m.m.EnqueueMessage(uid, message)
}

func (m *manager) AllClient() []int64 {
	return m.m.AllClient()
}
