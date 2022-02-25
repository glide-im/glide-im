package gateway

import (
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/message"
)

type manager struct {
	appId  int64
	m      *client.DefaultManager
	myAddr string
}

func newManager(etcd []string, myAddr string) (*manager, error) {
	ret := &manager{}
	ret.myAddr = myAddr
	ret.m = client.NewDefaultManager()

	var err error
	return ret, err
}

func (m *manager) AddClient(uid int64, cs client.IClient) {

}

func (m *manager) ClientConnected(conn conn.Connection) int64 {

	return 0
}

func (m *manager) ClientSignIn(oldUid int64, uid int64, device int64) {

	m.m.ClientSignIn(oldUid, uid, device)
}

func (m *manager) ClientLogout(uid int64, device int64) {

	m.m.ClientLogout(uid, device)
}

func (m *manager) EnqueueMessage(uid int64, device int64, message *message.Message) {
	m.m.EnqueueMessage(uid, device, message)
}
