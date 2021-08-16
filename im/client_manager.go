package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

var ClientManager = &clientManager{
	mutex:   NewMutex(),
	clients: map[int64]*Client{},
}

type clientManager struct {
	*mutex
	clients map[int64]*Client
}

func (c *clientManager) ClientSignIn(client *Client) {
	c.clients[client.uid] = client
}

func (c *clientManager) ClientSignOut(client *Client) {
	delete(c.clients, client.uid)
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

func (c *clientManager) EnqueueMessageMulti(uid int64, msg ...*entity.Message) {
	for _, message := range msg {
		c.EnqueueMessage(uid, message)
	}
}

func (c *clientManager) EnqueueMessage(uid int64, msg *entity.Message) {
	client := c.clients[uid]
	if c.IsOnline(uid) {
		if msg.Seq <= 0 {
			msg.Seq = client.getNextSeq()
		}
		client.EnqueueMessage(msg)
	} else {
		// TODO user offline
	}
}

func (c *clientManager) IsOnline(uid int64) bool {
	client, online := c.clients[uid]
	return online && !client.closed.Get()
}

func (c *clientManager) Update() {
	for _, client := range c.clients {
		if client.closed.Get() {
			c.ClientSignOut(client)
		}
	}
}

func (c *clientManager) AddGroup(uid int64, gid int64) {
	client, ok := c.clients[uid]
	if ok {
		client.AddGroup(gid)
	}
}

func (c *clientManager) RemoveGroup(uid int64, gid int64) {
	client, ok := c.clients[uid]
	if ok {
		client.RemoveGroup(gid)
	}
}

func (c *clientManager) AllClient() []int64 {
	var ret []int64
	for k := range c.clients {
		ret = append(ret, k)
	}
	return ret
}
