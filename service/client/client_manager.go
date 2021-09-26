package client

import (
	"fmt"
	"go_im/im"
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
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *manager) ClientConnected(conn conn.Connection) int64 {
	uid := m.m.ClientConnected(conn)
	uidTag := fmt.Sprintf("uid_%d", uid)
	err := m.router.SetTag("client", uidTag, m.myAddr)
	if err != nil {
		logger.E("set route tag error", err)
		return 0
	}
	return uid
}

func (m *manager) ClientSignIn(oldUid int64, uid int64, device int64) {
	err := m.router.RemoveTag("client", fmt.Sprintf("uid_%d", oldUid))
	if err != nil {

	}
	uidTag := fmt.Sprintf("uid_%d", uid)
	err = m.router.SetTag("client", uidTag, m.myAddr)
	if err != nil {

	}
	m.m.ClientSignIn(oldUid, uid, device)
}

func (m *manager) ClientLogout(uid int64) {
	err := m.router.RemoveTag("client", fmt.Sprintf("uid_%d", uid))
	if err != nil {

	}
	m.m.ClientLogout(uid)
}

func (m *manager) HandleMessage(from int64, message *message.Message) error {
	return m.m.HandleMessage(from, message)
}

func (m *manager) EnqueueMessage(uid int64, message *message.Message) {
	m.m.EnqueueMessage(uid, message)
}

func (m *manager) AllClient() []int64 {
	return m.m.AllClient()
}
