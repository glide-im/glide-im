package im

import (
	"errors"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/dao/uid"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
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

func (c *ClientManagerImpl) ClientSignIn(oldUid, uid_ int64, device int64) {
	cl := c.clients.Get(oldUid)
	if cl == nil {
		return
	}
	cl.SignIn(uid_, device)
	c.clients.Delete(oldUid)
	c.clients.Put(uid_, cl)
}

func (c *ClientManagerImpl) UserLogout(uid int64) {
	cl := c.clients.Get(uid)
	if cl == nil {
		return
	}
	c.clients.Delete(uid)
	cl.Exit()
}

func (c *ClientManagerImpl) HandleMessage(from int64, msg *message.Message) error {
	if msg.Action.Contains(message.ActionApi) {
		api.Handle(from, msg)
		return nil
	}
	switch msg.Action {
	case message.ActionChatMessage:
		return c.dispatchChatMessage(from, msg)
	case message.ActionGroupMessage:
		return group.Manager.DispatchMessage(from, msg)
	default:
		// unknown message type
	}
	return nil
}

func (c *ClientManagerImpl) dispatchChatMessage(from int64, msg *message.Message) error {
	senderMsg := new(client.SenderChatMessage)
	err := msg.DeserializeData(senderMsg)
	if err != nil {
		logger.E("sender chat senderMsg ", err)
		return err
	}
	logger.D("HandleMessage(from=%d): cid=%d, senderMsg=%s", from, senderMsg.Cid, senderMsg.Message)

	if senderMsg.Cid <= 0 {
		return errors.New("chat not create")
	}

	// update sender read time
	_ = dao.ChatDao.UpdateChatEnterTime(senderMsg.UcId)

	// insert message to chat
	chatMsg, err := dao.ChatDao.NewChatMessage(senderMsg.Cid, from, senderMsg.Message, senderMsg.MessageType)
	if err != nil {
		return err
	}
	affirm := message.NewMessage(msg.Seq, msg.Action, chatMsg)
	// send success, return chat message
	c.EnqueueMessage(from, affirm)

	return c.dispatch(from, chatMsg, senderMsg)
}

func (c *ClientManagerImpl) dispatch(from int64, chatMsg *dao.ChatMessage, senderMsg *client.SenderChatMessage) error {

	// update receiver's list chat
	uChat, err := dao.ChatDao.UpdateUserChatMsgTime(senderMsg.Cid, senderMsg.TargetId)
	if err != nil {
		return err
	}

	receiverMsg := client.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         senderMsg.Cid,
		UcId:        uChat.UcId,
		Sender:      from,
		MessageType: senderMsg.MessageType,
		Message:     senderMsg.Message,
		SendAt:      chatMsg.SendAt,
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(senderMsg.TargetId, dispatchMsg)

	return nil
}

func (c *ClientManagerImpl) EnqueueMessage(uid int64, msg *message.Message) {
	cl := c.clients.Get(uid)
	online := cl != nil && !cl.Closed()
	if online {
		//goland:noinspection GoNilness
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
	clients map[int64]*client.Client
}

func newClientMap() *clientMap {
	ret := new(clientMap)
	ret.Mutex = new(comm.Mutex)
	ret.clients = make(map[int64]*client.Client)
	return ret
}

func (g *clientMap) Size() int {
	return len(g.clients)
}

func (g *clientMap) Get(uid int64) *client.Client {
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

func (g *clientMap) Put(uid int64, client *client.Client) {
	defer g.LockUtilReturn()()
	g.clients[uid] = client
}

func (g *clientMap) Delete(uid int64) {
	defer g.LockUtilReturn()()
	delete(g.clients, uid)
}
