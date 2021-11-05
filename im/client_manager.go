package im

import (
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
)

type ClientManagerImpl struct {
	clients *clients
}

func NewClientManager() *ClientManagerImpl {
	ret := new(ClientManagerImpl)
	ret.clients = newClients()
	return ret
}

func (c *ClientManagerImpl) ClientConnected(conn conn.Connection) int64 {
	statistics.SConnEnter()

	// 获取一个临时 uid 标识这个连接
	connUid := uid.GenTemp()
	ret := client.NewClient(conn)
	ret.SetID(connUid, client.DeviceUnknown)
	c.clients.add(connUid, client.DeviceUnknown, ret)
	// 开始处理连接的消息
	ret.Run()
	return connUid
}

func (c *ClientManagerImpl) AddClient(uid int64, cs client.IClient) {
	c.clients.add(uid, client.DeviceUnknown, cs)
}

// ClientSignIn 客户端登录, id 为连接时使用的临时标识, uid 为用户标识, device 用于区分不同设备
func (c *ClientManagerImpl) ClientSignIn(id, uid int64, device int64) {
	ds := c.clients.get(id)
	if ds == nil || ds.size() == 0 {
		// 该客户端不存在
		logger.E("attempt to sign in a nonexistent client, id=%d", id)
		return
	}
	cli := ds.get(client.DeviceUnknown)
	cli.SetID(uid, device)

	// 移除临时 id 标识使用 uid 标记
	c.clients.delete(id, client.DeviceUnknown)

	loggedIn := c.clients.get(uid)
	if loggedIn != nil {
		log := loggedIn.get(device)
		if log != nil {
			// "Your account is logged in on another device"
			log.Exit(client.ExitCodeLoginMutex)
		}

		loggedIn.put(device, cli)
		//c.EnqueueMessage(uid, device, nil)
	} else {
		c.clients.add(uid, device, cli)
	}
}

func (c *ClientManagerImpl) ClientLogout(uid int64, device int64) {
	cl := c.clients.get(uid)
	if cl == nil || cl.size() == 0 {
		logger.E("uid is not sign in, uid=%d", uid)
		return
	}
	logDevice := cl.get(device)
	if logDevice == nil {
		logger.E("device not exist")
		return
	}
	logger.I("client logout, uid=%d, device=%d", uid, device)
	logDevice.Exit(client.ExitCodeByUser)
	cl.remove(device)
}

func (c *ClientManagerImpl) EnqueueMessage(uid int64, device int64, msg *message.Message) {
	ds := c.clients.get(uid)
	if ds == nil || ds.size() == 0 {
		// offline
		return
	}
	ds.foreach(func(deviceId int64, c client.IClient) {
		if device != 0 && deviceId != device {
			return
		}
		if c.Closed() {
			// TODO 2021-10-27 client is offline, store
		} else {
			c.EnqueueMessage(msg)
		}
	})
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

type devices struct {
	ds map[int64]client.IClient
}

func (d *devices) put(device int64, cli client.IClient) {
	d.ds[device] = cli
}

func (d *devices) get(device int64) client.IClient {
	return d.ds[device]
}

func (d *devices) remove(device int64) {
	delete(d.ds, device)
}

func (d *devices) foreach(f func(device int64, c client.IClient)) {
	for k, v := range d.ds {
		f(k, v)
	}
}
func (d *devices) size() int {
	return len(d.ds)
}

type clients struct {
	*comm.Mutex
	clients map[int64]*devices
}

func newClients() *clients {
	ret := new(clients)
	ret.Mutex = new(comm.Mutex)
	ret.clients = make(map[int64]*devices)
	return ret
}

func (g *clients) size() int {
	return len(g.clients)
}

func (g *clients) get(uid int64) *devices {
	defer g.LockUtilReturn()()
	cl, ok := g.clients[uid]
	if ok && cl.size() != 0 {
		return cl
	}
	return nil
}

func (g *clients) contains(uid int64) bool {
	_, ok := g.clients[uid]
	return ok
}

func (g *clients) add(uid int64, device int64, c client.IClient) {
	defer g.LockUtilReturn()()
	cs, ok := g.clients[uid]
	if ok {
		cs.put(device, c)
	} else {
		d := &devices{map[int64]client.IClient{}}
		d.put(device, c)
		g.clients[uid] = d
	}
}

func (g *clients) delete(uid int64, device int64) {
	defer g.LockUtilReturn()()
	d, ok := g.clients[uid]
	if ok {
		d.remove(device)
		if d.size() == 0 {
			delete(g.clients, uid)
		}
	}
}
