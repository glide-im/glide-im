package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

var ClientManager = &clientManager{clients: map[int64]*Client{}}

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
	uChat, err := dao.ChatDao.UpdateUserChatMsgTime(senderMsg.Cid, senderMsg.Receiver)
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

	dispatchMsg := entity.NewMessage(-1, entity.ActionChatMessage)

	if err = dispatchMsg.SetData(receiverMsg); err != nil {
		return err
	}

	if c.EnqueueMessage(senderMsg.Receiver, dispatchMsg) {
		// offline
	}

	return nil
}

func (c *clientManager) EnqueueMessage(uid int64, msg *entity.Message) bool {
	client, ok := c.clients[uid]
	if ok {
		if client.closed.Get() {
			ok = false
		} else {
			if msg.Seq == -1 {
				msg.Seq = client.getNextSeq()
			}
			if uid <= 0 {
				client.EnqueueMessage(entity.NewSimpleMessage(msg.Seq, entity.RespActionFailed, "unauthorized"))
				return false
			}
			client.EnqueueMessage(msg)
		}
	}
	return ok
}

func (c *clientManager) IsOnline(uid int64) bool {
	_, online := c.clients[uid]
	return online
}

func (c *clientManager) Update() {
	for _, client := range c.clients {
		if client.closed.Get() {
			c.ClientSignOut(client)
		}
	}
}

func (c *clientManager) AllClient() map[int64]*Client {
	return c.clients
}

func (c *clientManager) GetClient(uid int64) *Client {
	return c.clients[uid]
}
