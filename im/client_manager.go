package im

import (
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/dao/uid"
	"go_im/im/message"
)

type ClientManagerImpl struct {
	clients *clientMap
}

func NewClientManager() *ClientManagerImpl {
	ret := new(ClientManagerImpl)
	ret.clients = newClientMap()
	return ret
}

func (c *ClientManagerImpl) ClientConnected(conn conn.Connection) int64 {
	// 获取一个临时 uid 标识这个连接
	connUid := uid.GenTemp()
	ret := client.NewClient(conn, connUid)
	c.clients.Put(connUid, ret)
	// 开始处理连接的消息
	ret.Run()
	return connUid
}

func (c *ClientManagerImpl) AddClient(uid int64, cs client.IClient) {
	c.clients.Put(uid, cs)
}

func (c *ClientManagerImpl) ClientSignIn(oldUid, uid_ int64, device int64) {
	cl := c.clients.Get(oldUid)
	if cl == nil {
		return
	}
	cl.SetID(uid_)
	c.clients.Delete(oldUid)
	c.clients.Put(uid_, cl)
}

func (c *ClientManagerImpl) ClientLogout(uid int64) {
	cl := c.clients.Get(uid)
	if cl == nil {
		return
	}
	c.clients.Delete(uid)
	cl.Exit()
}

func (c *ClientManagerImpl) EnqueueMessage(uid int64, msg *message.Message) {
	cl := c.clients.Get(uid)
	online := cl != nil && !cl.Closed()
	if online {
		cl.EnqueueMessage(msg)
	} else {
		// TODO user offline
	}
}

func (c *ClientManagerImpl) AllClient() []int64 {
	var ret []int64
	for k := range c.clients.clients {
		if k > 0 {
			ret = append(ret, k)
		}
	}
	return ret
}

//////////////////////////////////////////////////////////////////////////////

type clientMap struct {
	*comm.Mutex
	clients map[int64]client.IClient
}

func newClientMap() *clientMap {
	ret := new(clientMap)
	ret.Mutex = new(comm.Mutex)
	ret.clients = make(map[int64]client.IClient)
	return ret
}

func (g *clientMap) Size() int {
	return len(g.clients)
}

func (g *clientMap) Get(uid int64) client.IClient {
	defer g.LockUtilReturn()()
	cl, ok := g.clients[uid]
	if ok {
		return cl
	}
	return nil
}

func (g *clientMap) Contains(uid int64) bool {
	_, ok := g.clients[uid]
	return ok
}

func (g *clientMap) Put(uid int64, client client.IClient) {
	defer g.LockUtilReturn()()
	g.clients[uid] = client
}

func (g *clientMap) Delete(uid int64) {
	defer g.LockUtilReturn()()
	delete(g.clients, uid)
}
