package im

import (
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

func (c *clientManager) SendChatMessage(from int64, message *entity.Message) error {

	msg := new(entity.ChatMessage)
	if err := message.DeserializeData(msg); err != nil {
		return err
	}

	if c.EnqueueMessage(msg.Target, message) {
		// offline
	}
	if err := dao.MessageDao.NewChatMessage(msg.ChatId, msg.Message, msg.MessageType); err != nil {
		return err
	}

	c.EnqueueMessage(from, entity.NewAckMessage(message.Seq))
	return nil
}

func (c *clientManager) EnqueueMessage(uid int64, msg *entity.Message) bool {
	client, ok := c.clients[uid]
	if ok {
		if client.closed {
			ok = false
		} else {
			client.EnqueueMessage(msg)
		}
	}
	return ok
}

func (c *clientManager) GetClient(uid int64) *Client {
	return c.clients[uid]
}
