package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

var ClientManager = newClientManager()

type clientManager struct {
	clients     *clientMap
	nextConnUid *AtomicInt64
}

func newClientManager() *clientManager {
	ret := new(clientManager)
	ret.clients = newClientMap()
	ret.nextConnUid = new(AtomicInt64)
	ret.nextConnUid.Set(-1)
	return ret
}

func (c *clientManager) ClientConnected(conn Connection) {
	client := NewClient(conn, c.nextConnUid.Get())
	c.nextConnUid.Set(client.uid - 1)
	c.clients.Put(client.uid, client)
	client.Run()
}

func (c *clientManager) ClientSignIn(oldUid, uid int64, device int64) {
	logger.D("ClientManager.ClientSignIn: connUid=%d, uid=%d", oldUid, uid)

	client := c.clients.Get(oldUid)
	if client == nil {
		return
	}
	client.uid = uid
	client.deviceId = device
	c.clients.Delete(oldUid)
	c.clients.Put(uid, client)
}

func (c *clientManager) UserLogout(uid int64) {
	logger.D("ClientManager.UserLogout: uid=%d", uid)
	c.clients.Delete(uid)
}

func (c *clientManager) DispatchMessage(from int64, message *entity.Message) error {

	senderMsg := new(entity.SenderChatMessage)
	if err := message.DeserializeData(senderMsg); err != nil {
		logger.E("sender chat senderMsg ", err)
		return err
	}
	logger.D("DispatchMessage(from=%d): cid=%d, senderMsg=%s", from, senderMsg.Cid, senderMsg.Message)

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
	affirm := entity.NewMessage(message.Seq, message.Action)
	if err = affirm.SetData(chatMsg); err != nil {
		return err
	}
	// send success, return chat message
	c.EnqueueMessage(from, affirm)

	// update receiver's list chat
	uChat, err := dao.ChatDao.UpdateUserChatMsgTime(senderMsg.Cid, senderMsg.TargetId)
	if err != nil {
		return err
	}

	receiverMsg := entity.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         senderMsg.Cid,
		UcId:        uChat.UcId,
		Sender:      from,
		MessageType: senderMsg.MessageType,
		Message:     senderMsg.Message,
		SendAt:      chatMsg.SendAt,
	}

	dispatchMsg := entity.NewMessage2(-1, entity.ActionChatMessage, receiverMsg)

	c.EnqueueMessage(senderMsg.TargetId, dispatchMsg)
	return nil
}

func (c *clientManager) EnqueueMessage(uid int64, msg *entity.Message) {
	client := c.clients.Get(uid)
	if c.IsOnline(uid) {
		client.EnqueueMessage(msg)
	} else {
		// TODO user offline
	}
}

func (c *clientManager) IsOnline(uid int64) bool {
	client := c.clients.Get(uid)
	return client != nil && client.uid > 0 && !client.closed.Get()
}

func (c *clientManager) Update() {
	for _, client := range c.clients.clients {
		if client.closed.Get() {
			c.UserLogout(client.uid)
		}
	}
}

func (c *clientManager) AllClient() []int64 {
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
	*mutex
	clients map[int64]*Client
}

func newClientMap() *clientMap {
	ret := new(clientMap)
	ret.mutex = new(mutex)
	ret.clients = make(map[int64]*Client)
	return ret
}

func (g *clientMap) Size() int {
	return len(g.clients)
}

func (g *clientMap) Get(uid int64) *Client {
	defer g.LockUtilReturn()()
	client, ok := g.clients[uid]
	if ok {
		return client
	}
	return nil
}

func (g *clientMap) Put(uid int64, client *Client) {
	defer g.LockUtilReturn()()
	g.clients[uid] = client
}

func (g *clientMap) Delete(uid int64) {
	defer g.LockUtilReturn()()
	delete(g.clients, uid)
}
