package im

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
	"runtime/debug"
)

// execPool 100 capacity goroutine pool, 假设每个消息处理需要10ms, 一个协程则每秒能处理100条消息
var execPool *ants.Pool

func init() {
	client.MessageHandleFunc = messageHandler

	var err error
	execPool, err = ants.NewPool(100,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(onHandleMessagePanic),
		//ants.WithPreAlloc(true),
	)
	if err != nil {
		panic(err)
	}
}

// messageHandler handle and dispatch client message
func messageHandler(from int64, device int64, msg *message.Message) {
	err := execPool.Submit(func() {
		switch msg.Action {
		case message.ActionChatMessage:
			dispatchChatMessage(from, msg)
		case message.ActionGroupMessage:
			dispatchGroupMsg(from, msg)
		case message.ActionCSMessage:
			dispatchCustomerServiceMsg(from, msg)
		default:
			if msg.Action.Contains(message.ActionApi) {
				api.Handle(from, msg)
			} else {
				client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotify, "unknown action"))
				logger.W("receive a unknown action message")
			}
		}
	})
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotify, "internal server error"))
		logger.E("async handle message error", err)
	}
}

func dispatchGroupMsg(from int64, msg *message.Message) {

	groupMsg := new(client.GroupMessage)
	err := msg.DeserializeData(groupMsg)
	if err != nil {
		logger.E("dispatch group message deserialize error", err)
		return
	}
	groupMsg.Sender = from
	group.Manager.DispatchMessage(groupMsg.TargetId, groupMsg)
}

func dispatchCustomerServiceMsg(from int64, msg *message.Message) {
	csMsg := new(client.CustomerServiceMessage)
	err := msg.DeserializeData(csMsg)
	csMsg.Sender = from
	if err != nil {
		logger.E("cs message", err)
		return
	}
	// 发送消息给客服
	client.EnqueueMessage(csMsg.CsId, msg)
}

func dispatchChatMessage(from int64, msg *message.Message) {
	senderMsg := new(client.SenderChatMessage)
	err := msg.DeserializeData(senderMsg)
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotify, "send message failed"))
		logger.E("sender chat senderMsg ", err)
		return
	}

	if senderMsg.Cid <= 0 {
		logger.E("dispatch message", "chat not create")
	}

	// update sender read time
	_ = dao.ChatDao.UpdateChatEnterTime(senderMsg.UcId)

	// insert message to chat
	chatMsg, err := dao.ChatDao.NewChatMessage(senderMsg.Cid, from, senderMsg.Message, senderMsg.MessageType)
	if err != nil {
		return
	}
	affirm := message.NewMessage(msg.Seq, msg.Action, chatMsg)
	// send success, return chat message
	client.EnqueueMessage(from, affirm)

	dispatch(from, chatMsg, senderMsg)
}

func dispatch(from int64, chatMsg *dao.ChatMessage, senderMsg *client.SenderChatMessage) {

	receiverMsg := client.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         senderMsg.Cid,
		Sender:      from,
		MessageType: senderMsg.MessageType,
		Message:     senderMsg.Message,
		SendAt:      chatMsg.SendAt,
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(senderMsg.TargetId, dispatchMsg)
}

func onHandleMessagePanic(i interface{}) {
	debug.PrintStack()
	logger.E("handler message panic", i)
}
