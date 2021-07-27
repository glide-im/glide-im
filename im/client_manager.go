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

	msg := new(entity.SenderChatMessage)
	if err := message.DeserializeData(msg); err != nil {
		logger.E("sender chat msg ", err)
		return err
	}
	logger.D("Chat message from=%d, cid=%d, msg=%s", from, msg.Cid, msg.Message)

	if msg.Cid <= 0 {
		return errors.New("chat not create")
	}

	// update sender read time
	_ = dao.MessageDao.UpdateChatEnterTime(msg.UcId)

	// insert message to chat
	chatMsg, err := dao.MessageDao.NewChatMessage(msg.Cid, from, msg.Message, msg.MessageType)
	if err != nil {
		return err
	}

	// update receiver's list chat
	uChat, err := dao.MessageDao.UpdateUserChatMsgTime(msg.Cid, msg.Receiver)
	if err != nil {
		return err
	}

	rm := entity.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         msg.Cid,
		UcId:        uChat.UcId,
		Sender:      from,
		MessageType: msg.MessageType,
		Message:     msg.Message,
		SendAt:      chatMsg.SendAt,
	}

	dispatchMsg := entity.NewMessage(-1, entity.ActionChatMessage)

	if err = dispatchMsg.SetData(rm); err != nil {
		return err
	}

	if c.EnqueueMessage(msg.Receiver, dispatchMsg) {
		// offline
	}

	c.EnqueueMessage(from, entity.NewAckMessage(message.Seq))
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

func (c *clientManager) AllClient() map[int64]*Client {
	return c.clients
}

func (c *clientManager) GetClient(uid int64) *Client {
	return c.clients[uid]
}
