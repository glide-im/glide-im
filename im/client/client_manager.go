package client

import (
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
	"strconv"
)

// Manager 客户端管理入口
var Manager IClientManager = NewDefaultManager()

type CommonInterface interface {
	// ClientSignIn 给一个已存在的客户端设置一个新的 id, 若 uid 已存在, 则新增一个 device 共享这个 id
	ClientSignIn(oldUid int64, uid int64, device int64)

	// ClientLogout 指定 uid, device 的客户端退出
	ClientLogout(uid int64, device int64)

	// EnqueueMessage 尝试将消息放入指定 uid 的客户端
	EnqueueMessage(uid int64, device int64, message *message.Message)

	// IsDeviceOnline 返回指定 uid 的用户设备是否在线
	IsDeviceOnline(uid, device int64) bool

	// IsOnline 返回指定 uid 的用户是否在线
	IsOnline(uid int64) bool

	// AllClient 返回所有的客户端 id
	AllClient() []int64
}

// IClientManager 管理所有客户端的创建, 消息派发, 退出等
type IClientManager interface {

	// ClientConnected 当一个用户连接建立后, 由该方法创建 IClient 实例 Client 并管理该连接, 返回该由连接创建客户端的标识 id
	// 返回的标识 id 是一个临时 id, 后续连接认证后会改变
	ClientConnected(conn conn.Connection) int64

	// AddClient 用于手段创建一个 IClient, 方便自定义临时 uid 以及其他的 IClient 实现
	AddClient(uid int64, cs IClient)

	CommonInterface
}

// EnqueueMessage Manager.EnqueueMessage 的快捷方法, 预留一个位置对消息入队列进行一些预处理
func EnqueueMessage(uid int64, message *message.Message) {
	//
	Manager.EnqueueMessage(uid, 0, message)
}

func EnqueueMessageToDevice(uid int64, device int64, message *message.Message) {
	Manager.EnqueueMessage(uid, device, message)
}

type DefaultManager struct {
	clients *clients
}

func NewDefaultManager() *DefaultManager {
	ret := new(DefaultManager)
	ret.clients = newClients()
	return ret
}

func (c *DefaultManager) ClientConnected(conn conn.Connection) int64 {
	statistics.SConnEnter()

	// 获取一个临时 uid 标识这个连接
	connUid := uid.GenTemp()
	ret := newClient(conn)
	ret.SetID(connUid, 0)
	c.clients.add(connUid, 0, ret)
	// 开始处理连接的消息
	ret.Run()
	return connUid
}

func (c *DefaultManager) AddClient(uid int64, cs IClient) {
	c.clients.add(uid, 0, cs)
}

// ClientSignIn 客户端登录, id 为连接时使用的临时标识, uid 为用户标识, device 用于区分不同设备
func (c *DefaultManager) ClientSignIn(id, uid_ int64, device int64) {
	logger.D("client sign in origin-id=%d, uid=%d", id, uid_)
	tempDs := c.clients.get(id)
	if tempDs == nil || tempDs.size() == 0 {
		// 该客户端不存在
		logger.W("attempt to sign in a nonexistent client, id=%d", id)
		return
	}
	client := tempDs.get(0)
	logged := c.clients.get(uid_)
	if logged != nil && logged.size() > 0 {
		// 多设备登录
		existing := logged.get(device)
		if existing != nil {
			logger.D("multi device login mutex, uid=%d, device=%d", uid_, device)
			existing.SetID(uid.GenTemp(), 0)
			// "Your account is logged in on another device"
			existing.EnqueueMessage(message.NewMessage(0, message.ActionKickOut, "Your account is logged in on another device"))
			existing.Exit()
			logged.remove(device)
		}
		if logged.size() > 0 {
			EnqueueMessage(uid_, message.NewMessage(0, message.ActionNotify, "multi device login, device="+strconv.FormatInt(device, 10)))
		}
		logged.put(device, client)
	} else {
		// 单设备登录
		c.clients.add(uid_, device, client)
	}
	client.SetID(uid_, device)
	// 删除临时 id
	c.clients.delete(id, 0)
}

func (c *DefaultManager) ClientLogout(uid int64, device int64) {
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
	logDevice.Exit()
	cl.remove(device)
}

func (c *DefaultManager) EnqueueMessage(uid int64, device int64, msg *message.Message) {
	ds := c.clients.get(uid)
	if ds == nil || ds.size() == 0 {
		// offline
		return
	}
	ds.foreach(func(deviceId int64, c IClient) {
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

func (c *DefaultManager) IsOnline(uid int64) bool {
	ds := c.clients.get(uid)
	if ds == nil {
		return false
	}
	return ds.size() > 0
}

func (c *DefaultManager) IsDeviceOnline(uid, device int64) bool {
	ds := c.clients.get(uid)
	if ds == nil {
		return false
	}
	return ds.get(device) != nil
}

func (c *DefaultManager) AllClient() []int64 {
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
	ds map[int64]IClient
}

func (d *devices) put(device int64, cli IClient) {
	d.ds[device] = cli
}

func (d *devices) get(device int64) IClient {
	return d.ds[device]
}

func (d *devices) remove(device int64) {
	delete(d.ds, device)
}

func (d *devices) foreach(f func(device int64, c IClient)) {
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

func (g *clients) add(uid int64, device int64, c IClient) {
	defer g.LockUtilReturn()()
	cs, ok := g.clients[uid]
	if ok {
		cs.put(device, c)
	} else {
		d := &devices{map[int64]IClient{}}
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
