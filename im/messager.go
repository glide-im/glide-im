package im

import (
	"errors"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

func init() {
	client.MessageHandleFunc = messageHandler
}

// messageHandler handle and dispatch client message
func messageHandler(from int64, msg *message.Message) error {
	switch msg.Action {
	case message.ActionChatMessage:
		return dispatchChatMessage(from, msg)
	case message.ActionGroupMessage:
		return group.Manager.DispatchMessage(from, msg)
	case message.ActionCSMessage:
		return dispatchCustomerServiceMsg(from, msg)
	default:
		if msg.Action.Contains(message.ActionApi) {
			api.Handle(from, msg)
		} else {
			// unknown type
		}
	}
	return nil
}

func dispatchCustomerServiceMsg(from int64, msg *message.Message) error {
	csMsg := new(client.CustomerServiceMessage)
	err := msg.DeserializeData(csMsg)
	csMsg.From = from
	if err != nil {
		logger.E("cs message", err)
		return err
	}
	client.EnqueueMessage(csMsg.CsId, msg)
	return nil
}

func dispatchChatMessage(from int64, msg *message.Message) error {
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
	client.EnqueueMessage(from, affirm)

	return dispatch(from, chatMsg, senderMsg)
}

func dispatch(from int64, chatMsg *dao.ChatMessage, senderMsg *client.SenderChatMessage) error {

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
