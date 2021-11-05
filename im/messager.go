package im

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
)

// execPool 100 capacity goroutine pool, 假设每个消息处理需要10ms, 一个协程则每秒能处理100条消息
var execPool *ants.Pool

func init() {
	client.MessageHandleFunc = messageHandler

	var err error
	execPool, err = ants.NewPool(200000,
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
		statistics.SMsgInput()
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
		logger.E("async handle message error %v", err)
	}
}

func dispatchGroupMsg(from int64, msg *message.Message) {
	groupMsg := new(message.GroupMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.Sender = from
	group.Manager.DispatchMessage(groupMsg.TargetId, groupMsg)
}

func dispatchCustomerServiceMsg(from int64, msg *message.Message) {
	csMsg := new(message.CustomerServiceMessage)
	if !unwrap(from, msg, csMsg) {
		return
	}
	// 发送消息给客服
	client.EnqueueMessage(csMsg.CsId, msg)
}

func dispatchChatMessage(from int64, msg *message.Message) {
	senderMsg := new(message.SenderChatMessage)
	if !unwrap(from, msg, senderMsg) {
		return
	}

	if senderMsg.Cid <= 0 {
		logger.E("chat not create, from=%d, to=%d", from, senderMsg.TargetId)
	}

	// update sender read time
	_ = dao.ChatDao.UpdateChatEnterTime(senderMsg.UcId)

	// insert message to chat
	chatMsg, err := dao.ChatDao.NewChatMessage(senderMsg.Cid, from, senderMsg.Message, senderMsg.MessageType)
	if err != nil {
		return
	}
	ackMessage(msg.Seq, chatMsg)
	dispatch(from, chatMsg, senderMsg)
}

func dispatch(from int64, chatMsg *dao.ChatMessage, senderMsg *message.SenderChatMessage) {

	receiverMsg := message.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         senderMsg.Cid,
		Sender:      from,
		MessageType: senderMsg.MessageType,
		Message:     senderMsg.Message,
		SendAt:      chatMsg.SendAt.Unix(),
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(senderMsg.TargetId, dispatchMsg)
}

func onHandleMessagePanic(i interface{}) {
	logger.E("handler message panic, %v", i)
}

func ackMessage(seq int64, m *dao.ChatMessage) {
	ack := message.ChatMessageAck{
		Seq: seq,
		Mid: m.Mid,
	}
	ackResp := message.NewMessage(seq, message.ActionMessageAck, ack)
	client.EnqueueMessage(m.Sender, ackResp)
}

func unwrap(from int64, msg *message.Message, to interface{}) bool {
	err := msg.DeserializeData(to)
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(msg.Seq, message.ActionNotify, "send message failed"))
		logger.E("sender chat senderMsg %v", err)
		return false
	}
	return true
}
